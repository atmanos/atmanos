package runtime

import "unsafe"

func osinit() {
	ncpu = 1

	atmaninit()
}

func sigpanic() {}

func crash() {
	*(*int32)(nil) = 0
}

func goenvs() {}

//go:nowritebarrier
func newosproc(mp *m, stk unsafe.Pointer) {
	mp.tls[0] = uintptr(mp.id) // so 386 asm can find it
	if true {
		print("newosproc stk=", stk, " m=", mp, " g=", mp.g0, " id=", mp.id, "/", mp.tls[0], " ostk=", &mp, "\n")
	}

	taskcreate(
		unsafe.Pointer(mp),
		unsafe.Pointer(mp.g0),
		unsafe.Pointer(funcPC(mstart)),
		stk,
	)

	taskyield()
}

func resetcpuprofiler(hz int32) {}

func minit() {
	println("minit()")
}

//go:nosplit
func unminit() {
	println("unminit()")
}

//go:nosplit
func mpreinit(mp *m) {
	print("mpreinit(", unsafe.Pointer(mp), ")", "\n")
}

//go:nosplit
func msigsave(mp *m) {
	print("msigsave(", unsafe.Pointer(mp), ")", "\n")
}

//go:nosplit
func msigrestore(mp *m) {}

//go:nosplit
func sigblock() {}

//go:nosplit
func osyield() {
	taskyield()
}

// Create a semaphore, which will be assigned to m->waitsema.
// The zero value is treated as absence of any semaphore,
// so be sure to return a non-zero value.
//go:nosplit
func semacreate() uintptr {
	return 1
}

// If ns < 0, acquire m->waitsema and return 0.
// If ns >= 0, try to acquire m->waitsema for at most ns nanoseconds.
// Return 0 if the semaphore was acquired, -1 if interrupted or timed out.
//go:nosplit
//go:nowritebarrier
func semasleep(ns int64) int32 {
	var ret int32

	systemstack(func() {
		ret = semasleepInternal(ns)
	})

	return ret
}

//go:nowritebarrier
func semasleepInternal(ns int64) int32 {
	_g_ := getg()

	var addr = &_g_.m.waitsemacount

	if cas(addr, 1, 0) {
		return 0
	}

	var (
		waiter = &taskcurrent.semawaiter
		s      = &sleeptable[sleeptablekey(addr)]
	)

	atomicstorep1(unsafe.Pointer(&waiter.addr), unsafe.Pointer(addr))
	waiter.up = false
	s.add(waiter)

	tasksleep(ns)

	if !waiter.up {
		// interrupted or timed out
		s.remove(waiter)
		return -1
	}

	return 0
}

// Wake up mp, which is or will soon be sleeping on mp->waitsema.
//go:nosplit
//go:nowritebarrier
func semawakeup(mp *m) {
	var (
		addr = &mp.waitsemacount
		s    = &sleeptable[sleeptablekey(addr)]
	)

	waiter := s.removeWaiterOn(addr)
	if waiter == nil {
		xchg(addr, 1)
	} else {
		waiter.up = true
		taskwake(waiter.task)
	}
}

var (
	sleeptable [512]sema
)

func sleeptablekey(addr *uint32) int {
	a := uintptr(unsafe.Pointer(addr))
	return int(a & 511) // TODO: hash address
}

type sema struct {
	head, tail *semawaiter
}

func (s *sema) removeWaiterOn(addr *uint32) *semawaiter {
	for w := s.head; w != nil; w = w.next {
		if w.addr != addr {
			continue
		}

		s.remove(w)
		return w
	}

	return nil
}

//go:nowritebarrier
func (s *sema) remove(w *semawaiter) {
	if w.prev != nil {
		atomicstorep1(unsafe.Pointer(&w.prev.next), unsafe.Pointer(w.next))
	} else {
		atomicstorep1(unsafe.Pointer(&s.head), unsafe.Pointer(w.next))
	}

	if w.next != nil {
		atomicstorep1(unsafe.Pointer(&w.next.prev), unsafe.Pointer(w.prev))
	} else {
		atomicstorep1(unsafe.Pointer(&s.tail), unsafe.Pointer(w.prev))
	}
}

//go:nowritebarrier
func (s *sema) add(w *semawaiter) {
	if s.tail != nil {
		atomicstorep1(unsafe.Pointer(&s.tail.next), unsafe.Pointer(w))
		atomicstorep1(unsafe.Pointer(&w.prev), unsafe.Pointer(s.tail))
	} else {
		atomicstorep1(unsafe.Pointer(&s.head), unsafe.Pointer(w))
		atomicstorep1(unsafe.Pointer(&w.prev), nil)
	}

	atomicstorep1(unsafe.Pointer(&s.tail), unsafe.Pointer(w))
	w.next = nil
}

type semawaiter struct {
	addr *uint32
	task *Task
	up   bool

	next, prev *semawaiter
}
