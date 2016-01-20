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
	println("osyield()")
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
func semasleep(ns int64) int32 {
	print("semasleep(", ns, ")", "\n")
	_g_ := getg()
	addr := &_g_.m.waitsemacount

	var (
		waiter = &semawaiter{addr: addr, task: taskcurrent, up: false}
		s      = &sleeptable[sleeptablekey(addr)]
	)

	s.lock()

	if *addr > 0 {
		xadd(addr, -1)
		s.unlock()
		return 0
	}

	s.add(waiter)

	for !waiter.up {
		s.unlock()
		ns = tasksleep(ns)
		s.lock()

		if ns == 0 {
			break
		}
	}

	if !waiter.up {
		s.remove(waiter)
	}

	s.unlock()

	if waiter.up {
		return 0
	}

	return -1
}

// Wake up mp, which is or will soon be sleeping on mp->waitsema.
//go:nosplit
func semawakeup(mp *m) {
	var (
		addr = &mp.waitsemacount
		s    = &sleeptable[sleeptablekey(addr)]
	)

	s.lock()

	waiter := s.removeWaiterOn(addr)
	if waiter == nil {
		xadd(addr, 1)
	} else {
		waiter.up = true
		taskwake(waiter.task)
	}

	s.unlock()
}

var (
	sleeptable [512]sema
	sleeplocks [512]qlock
)

func sleeptablekey(addr *uint32) int {
	a := uintptr(unsafe.Pointer(addr))
	return int(a & 511) // TODO: hash address
}

type qlock struct {
	owner   *Task
	waiting TaskList
}

func (l *qlock) lock() {
	if l.owner == nil {
		l.owner = taskcurrent
		return
	}

	l.waiting.Add(taskcurrent)
	taskswitch()
}

func (l *qlock) unlock() {
	ready := l.waiting.Head
	l.owner = ready

	if ready != nil {
		l.waiting.Remove(ready)
		taskready(ready)
	}
}

type sema struct {
	*qlock

	head, tail *semawaiter
}

func (s *sema) removeWaiterOn(addr *uint32) *semawaiter {
	return nil
}

func (s *sema) remove(w *semawaiter) {
}

func (s *sema) add(w *semawaiter) {
	if s.tail != nil {
		s.tail.next = w
		w.prev = s.tail
	} else {
		s.head = w
		w.prev = nil
	}

	s.tail = w
	w.next = nil
}

type semawaiter struct {
	addr *uint32
	task *Task
	up   bool

	next, prev *semawaiter
}

func init() {
	for i := 0; i < 512; i++ {
		sleeptable[i].qlock = &sleeplocks[i]
	}
}
