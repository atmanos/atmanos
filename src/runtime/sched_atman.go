package runtime

import "unsafe"

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
)

type Task struct {
	ID    int
	State [256]byte

	Context Context

	WakeAt int64

	Next, Prev *Task
}

// taskcreate spawns a new task,
// executing the argumentless function fn
// on the provided stack stk
func taskcreate(mp, g0, fn, stk unsafe.Pointer) {
	// create hole in stack
	stk = unsafe.Pointer(uintptr(stk) - 64)

	// reserve stack space to create Task.
	taskSize := unsafe.Sizeof(Task{})
	stk = unsafe.Pointer(uintptr(stk) - taskSize)
	memclr(stk, taskSize)
	t := (*Task)(unsafe.Pointer(stk))

	// create hole in stack
	stk = unsafe.Pointer(uintptr(stk) - 64)

	stk = unsafe.Pointer(uintptr(stk) - 8)
	*(*uintptr)(stk) = uintptr(g0)

	stk = unsafe.Pointer(uintptr(stk) - 8)
	*(*uintptr)(stk) = uintptr(mp)

	// reserve 8 bytes of space for return value,
	// for calling compatibility with contextsave.
	stk = unsafe.Pointer(uintptr(stk) - 8)

	stk = unsafe.Pointer(uintptr(stk) - 8)
	*(*uintptr)(stk) = uintptr(fn)

	// reserve 8 bytes of space for return address,
	// to be filled in by contextload.
	stk = unsafe.Pointer(uintptr(stk) - 8)

	t.ID = taskid
	t.Context = Context{
		rsp: uintptr(stk),
		rip: funcPC(taskstart),
	}

	taskid++
	taskn++

	taskready(t)
}

//go:nosplit
func taskstart(fn, _, mp, gp unsafe.Pointer)

func taskready(t *Task) {
	t.WakeAt = 0
	taskrunqueue.Add(t)
}

func taskyield() {
	taskready(taskcurrent)
	taskswitch()
}

func taskswitch() {
	taskrunqueue.debug()

	var (
		taskprev = taskcurrent
		tasknext *Task
	)

	for {
		taskwakeready(nanotime())

		if tasknext = taskrunqueue.Head; tasknext != nil {
			break
		}

		if tasksleepqueue.Head == nil || tasksleepqueue.Head.WakeAt == -1 {
			panic("No runnable or timed sleep tasks to run")
		}

		osyield()
	}

	taskcurrent = tasknext
	taskrunqueue.Remove(taskcurrent)

	println("switching from", taskprev.ID, "to", taskcurrent.ID)
	contextswitch(&taskprev.Context, &taskcurrent.Context)
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

	if rem = sleepend - sleepstart; rem < 0 {
		rem = 0
	}

	return rem
}

// taskwake moves task from the sleep to the run queue.
func taskwake(task *Task) {
	tasksleepqueue.Remove(task)
	taskready(task)
}

func taskwakeready(at int64) {
	for {
		task := tasksleepqueue.Head
		if task == nil || task.WakeAt < 0 || task.WakeAt > at {
			return
		}

		tasksleepqueue.Remove(task)
		taskready(task)
	}
}

func taskexit() {
	panic("taskexit()")
}

type TaskList struct {
	Head, Tail *Task
}

func (l *TaskList) debug() {
	print("[")
	for t := l.Head; t != nil; t = t.Next {
		print("Task{ID: ", t.ID, "}, ")
	}
	print("]\n")
}

func (l *TaskList) Add(t *Task) {
	if l.Tail != nil {
		l.Tail.Next = t
		t.Prev = l.Tail
	} else {
		l.Head = t
		t.Prev = nil
	}

	l.Tail = t
	t.Next = nil
}

func (l *TaskList) Remove(t *Task) {
	if t.Prev != nil {
		t.Prev.Next = t.Next
	} else {
		l.Head = t.Next
	}

	if t.Next != nil {
		t.Next.Prev = t.Prev
	} else {
		l.Tail = t.Prev
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

// Context describes the state of a task
// for saving or restoring a task's execution context.
type Context cpuRegisters

func (c *Context) debug() {
	print(
		"Context{",
		"rsp=", unsafe.Pointer(c.rsp),
		" rip=", unsafe.Pointer(c.rip),
		" fs=", unsafe.Pointer(c.fs),
		"}", "\n",
	)
}

func contextswitch(from, to *Context) {
	if contextsave(from) == 0 {
		from.debug()
		to.debug()
		contextload(to)
	}
}

func contextsave(*Context) int
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
	es     uintptr
	ds     uintptr
	fs     uintptr
	gs     uintptr
}
