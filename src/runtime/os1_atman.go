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

	for {
		if _g_.m.waitsemacount > 0 {
			_g_.m.waitsemacount--
			return 0 // semaphore acquired
		}

		a := uintptr(unsafe.Pointer(&_g_.m.waitsemacount))
		key := int(a & 511) // TODO: hash address
		s := &sleeptable[key]

		s.qlock.lock()
		s.waiting.Add(taskcurrent)
		s.qlock.unlock()
		taskswitch()
	}
}

// Wake up mp, which is or will soon be sleeping on mp->waitsema.
//go:nosplit
func semawakeup(mp *m) int32 {
	print("semawakeup(", unsafe.Pointer(mp), ")", "\n")
	crash()
	return 0
}

var (
	sleeptable [512]sema
	sleeplocks [512]qlock
)

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
	qlock   *qlock
	waiting TaskList
}

func init() {
	for i := 0; i < 512; i++ {
		sleeptable[i].qlock = &sleeplocks[i]
	}
}
