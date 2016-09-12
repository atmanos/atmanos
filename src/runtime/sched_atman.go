package runtime

import (
	"runtime/internal/atomic"
	"unsafe"
)

// Atman implements a round-robin cooperative scheduler,
// heavily inspired by libtask.

var (
	taskid = 1
	taskn  = 1

	// taskcurrent is initialized as an empty Task with ID 0.
	// The first time a context switch is requested,
	// its context will be filled in and it will become
	// a normal Task.
	taskcurrent = &Task{ID: 0}

	taskrunqueue   TaskList
	tasksleepqueue TaskList
	taskcache      taskCache
)

func init() {
	taskcurrent.semawaiter.task = taskcurrent
}

type Task struct {
	ID int

	Context Context

	WakeAt int64

	Next, Prev *Task

	semawaiter
}

func (t *Task) debug() {
	println("Task{ID: ", t.ID, ", WakeAt: ", t.WakeAt, "}")
}

// taskcreate spawns a new task,
// executing the argumentless function fn
// on the provided stack stk
func taskcreate(mp, g0, fn, stk unsafe.Pointer) {
	// create hole in stack
	stk = unsafe.Pointer(uintptr(stk) - 64)

	stk = unsafe.Pointer(uintptr(stk) - 8)
	*(*uintptr)(stk) = uintptr(g0)

	stk = unsafe.Pointer(uintptr(stk) - 8)
	*(*uintptr)(stk) = uintptr(mp)

	stk = unsafe.Pointer(uintptr(stk) - 8)
	*(*uintptr)(stk) = uintptr(fn)

	t := taskcache.alloc()
	t.ID = taskid

	contextsave(&t.Context, 0)
	t.Context.r.rsp = uintptr(stk)
	t.Context.r.rip = funcPC(taskstart)
	atomic.StorepNoWB(unsafe.Pointer(&t.semawaiter.task), unsafe.Pointer(t))

	taskid++
	taskn++

	taskready(t)
}

//go:nosplit
func taskstart(fn, mp, gp unsafe.Pointer)

func taskready(t *Task) {
	t.WakeAt = 0
	taskrunqueue.Add(t)
}

//gc:nowritebarrier
func taskyield() {
	taskready(taskcurrent)
	taskswitch()
}

//go:nosplit
func taskswitch() {
	contextsave(&taskcurrent.Context, funcPC(taskschedule))
}

func taskschedule() {
	var tasknext *Task

	saved := irqDisable()
	for {
		taskwakeready(nanotime())

		if tasknext = taskrunqueue.Head; tasknext != nil {
			break
		}

		if tasksleepqueue.Head == nil || tasksleepqueue.Head.WakeAt == -1 {
			panic("No runnable or timed sleep tasks to run")
		}

		HYPERVISOR_set_timer_op(tasksleepqueue.Head.WakeAt)
		HYPERVISOR_sched_op(1, nil) // block
		irqDisable()
	}
	irqRestore(saved)

	atomic.StorepNoWB(unsafe.Pointer(&taskcurrent), unsafe.Pointer(tasknext))
	taskrunqueue.Remove(taskcurrent)
	contextload(&taskcurrent.Context)
}

func tasksleepus(us uint32) {
	ns := int64(us) * 1000

	for ns > 0 {
		ns = tasksleep(ns)
	}
}

// tasksleep puts the current task to sleep for up to ns.
// It returns the remaining sleep time if woken early.
// If ns is -1, rem will always be -1.
func tasksleep(ns int64) (rem int64) {
	sleepstart := nanotime()

	if ns == -1 {
		taskcurrent.WakeAt = -1
	} else {
		taskcurrent.WakeAt = sleepstart + ns
	}

	tasksleepqueue.AddByWakeAt(taskcurrent)
	taskswitch()

	sleepend := nanotime()

	if ns < 0 {
		return -1
	}

	if rem = ns - (sleepend - sleepstart); rem < 0 {
		rem = 0
	}

	return rem
}

// taskwake moves task from the sleep to the run queue.
//
// If task is already marked as awake, it does nothing.
func taskwake(task *Task) {
	if task.WakeAt != 0 {
		tasksleepqueue.Remove(task)
		taskready(task)
	}
}

func taskwakeready(at int64) {
	for {
		task := tasksleepqueue.Head
		if task == nil || task.WakeAt < 0 || task.WakeAt > at {
			return
		}
		taskwake(task)
	}
}

func taskexit() {
	throw("taskexit()")
	panic("taskexit()")
}

type TaskList struct {
	Head, Tail *Task
}

func (l *TaskList) debug() {
	println("[")
	for t := l.Head; t != nil; t = t.Next {
		print("  ")
		t.debug()
	}
	println("]")
}

//gc:nowritebarrier
func (l *TaskList) Add(t *Task) {
	if l.Tail != nil {
		atomic.StorepNoWB(unsafe.Pointer(&l.Tail.Next), unsafe.Pointer(t))
		atomic.StorepNoWB(unsafe.Pointer(&t.Prev), unsafe.Pointer(l.Tail))
	} else {
		atomic.StorepNoWB(unsafe.Pointer(&l.Head), unsafe.Pointer(t))
		t.Prev = nil
	}

	atomic.StorepNoWB(unsafe.Pointer(&l.Tail), unsafe.Pointer(t))
	t.Next = nil
}

//gc:nowritebarrier
func (l *TaskList) Remove(t *Task) {
	if t.Prev != nil {
		atomic.StorepNoWB(unsafe.Pointer(&t.Prev.Next), unsafe.Pointer(t.Next))
		// t.Prev.Next = t.Next
	} else {
		atomic.StorepNoWB(unsafe.Pointer(&l.Head), unsafe.Pointer(t.Next))
		// l.Head = t.Next
	}

	if t.Next != nil {
		atomic.StorepNoWB(unsafe.Pointer(&t.Next.Prev), unsafe.Pointer(t.Prev))
	} else {
		atomic.StorepNoWB(unsafe.Pointer(&l.Tail), unsafe.Pointer(t.Prev))
	}
}

func (l *TaskList) AddByWakeAt(t *Task) {
	if t.WakeAt < 0 {
		l.Add(t)
		return
	}

	for i := l.Head; i != nil; i = i.Next {
		if t.WakeAt > i.WakeAt && i.WakeAt >= 0 {
			continue
		}

		if i.Prev == nil {
			l.Head = t
		} else {
			i.Prev.Next = t
		}

		t.Prev = i.Prev
		t.Next = i
		i.Prev = t

		return
	}

	// no match, add to tail
	l.Add(t)
}

type taskCache struct {
	list TaskList
}

func (c *taskCache) alloc() *Task {
	if c.list.Head == nil {
		const taskSize = unsafe.Sizeof(Task{})
		n := _PAGESIZE / taskSize
		if n == 0 {
			n = 1
		}

		// Must be non-GC memory because the GC and runtime
		// depend on tasks.
		mem := persistentalloc(n*taskSize, 0, &memstats.other_sys)
		for i := uintptr(0); i < n; i++ {
			task := (*Task)(add(mem, i*taskSize))
			c.list.Add(task)
		}
	}

	task := c.list.Head
	c.list.Remove(task)
	return task
}

// Context describes the state of a task
// for saving or restoring a task's execution context.
type Context struct {
	r   cpuRegisters
	tls uintptr
}

func (c *Context) debug() {
	print(
		"Context{",
		"rsp=", unsafe.Pointer(c.r.rsp),
		" rip=", unsafe.Pointer(c.r.rip),
		" tls=", unsafe.Pointer(c.tls),
		"}", "\n",
	)
}

//go:nosplit
func contextsave(*Context, uintptr)

//go:nosplit
func contextload(*Context)

// cpuRegisters describes the state of a CPU
// on entry of a trap or interrupt.
type cpuRegisters struct {
	r15    uintptr
	r14    uintptr
	r13    uintptr
	r12    uintptr
	rbp    uintptr
	rbx    uintptr
	r11    uintptr
	r10    uintptr
	r9     uintptr
	r8     uintptr
	rax    uintptr
	rcx    uintptr
	rdx    uintptr
	rsi    uintptr
	rdi    uintptr
	code   uintptr // syscall number, error code, or IRQ number
	rip    uintptr
	cs     uintptr
	rflags uintptr
	rsp    uintptr
	ss     uintptr
}

func (r *cpuRegisters) debug() {
	print("cpuRegisters<", unsafe.Pointer(r), "> {")
	print(" r15=", unsafe.Pointer(r.r15))
	print(" r14=", unsafe.Pointer(r.r14))
	print(" r13=", unsafe.Pointer(r.r13))
	print(" r12=", unsafe.Pointer(r.r12))
	print(" rbp=", unsafe.Pointer(r.rbp))
	print(" rbx=", unsafe.Pointer(r.rbx))
	print(" r11=", unsafe.Pointer(r.r11))
	print(" r10=", unsafe.Pointer(r.r10))
	print(" r9=", unsafe.Pointer(r.r9))
	print(" r8=", unsafe.Pointer(r.r8))
	print(" rax=", unsafe.Pointer(r.rax))
	print(" rcx=", unsafe.Pointer(r.rcx))
	print(" rdx=", unsafe.Pointer(r.rdx))
	print(" rsi=", unsafe.Pointer(r.rsi))
	print(" rdi=", unsafe.Pointer(r.rdi))
	print(" code=", unsafe.Pointer(r.code))
	print(" rip=", unsafe.Pointer(r.rip))
	print(" cs=", unsafe.Pointer(r.cs))
	print(" rflags=", unsafe.Pointer(r.rflags))
	print(" rsp=", unsafe.Pointer(r.rsp))
	print(" ss=", unsafe.Pointer(r.ss))
	print(" }\n")
}
