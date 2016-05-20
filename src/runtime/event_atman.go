package runtime

import (
	"runtime/internal/atomic"
	"unsafe"
)

type eventHandler func(port uint32, r *cpuRegisters)

var eventHandlers [4096]eventHandler

func initEvents() {
	HYPERVISOR_set_callbacks(
		funcPC(eventCallbackASM),
		funcPC(eventFailsafe),
	)

	enableIRQ()
	bindVIRQ()
}

// clearbit atomically clears the bit at offset n from addr.
//
//go:nosplit
func clearbit(addr *uint64, n uint32)

//go:nosplit
func eventChanSend(port uint32) {
	op := struct{ port uint32 }{port: port}

	ret := HYPERVISOR_event_channel_op(
		4, // EVTCHNOP_send
		unsafe.Pointer(&op),
	)

	if ret != 0 {
		println("HYPERVISOR_event_channel_op returned", ret)
	}
}

//go:nosplit
func unmaskEventChan(port uint32) {
	info := _atman_shared_info
	clearbit(&info.EvtchnMask[0], port)
}

//go:nosplit
func clearEventChan(port uint32) {
	clearbit(&_atman_shared_info.EvtchnPending[0], port)
}

// enableIRQ enables interrupts clearing the event mask
// on the shared VCPU structure.
//
//go:nosplit
func enableIRQ() {
	irqRestore(0)
}

//go:nosplit
func irqDisable() uint8 {
	vcpu := &_atman_shared_info.VCPUInfo[0]

	mask := vcpu.UpcallMask
	vcpu.UpcallMask = 1
	return mask
}

//go:nosplit
func irqRestore(mask uint8) {
	vcpu := &_atman_shared_info.VCPUInfo[0]
	vcpu.UpcallMask = mask
}

func bindVIRQ() {
	op := struct {
		VIRQ uint32
		VCPU uint32

		// Port is set on success
		Port uint32
	}{
		VIRQ: 0, // VIRQ_TIMER,
		VCPU: 0,
	}

	ret := HYPERVISOR_event_channel_op(
		1, // EVTCHNOP_bind_virq
		unsafe.Pointer(&op),
	)

	if ret != 0 {
		println("HYPERVISOR_event_channel_op returned", ret)
	}

	bindEventHandler(op.Port, virqHandler)

	println("Bound VIRQ to Port", op.Port)
}

func virqHandler(port uint32, r *cpuRegisters) {
	taskwakeready(nanotime())
}

// bindEventHandler binds the handler f to events on port.
func bindEventHandler(port uint32, f eventHandler) {
	eventHandlers[port] = f
}

// handleEvents fires event handlers for any pending events.
//
//go:nosplit
func handleEvents(r *cpuRegisters) {
	_atman_shared_info.VCPUInfo[0].UpcallPending = 0

	sel := atomic.Xchg64(&_atman_shared_info.VCPUInfo[0].PendingSel, 0)

	for i := uint32(0); i < 64; i++ {
		// each set bit in sel is an index into EvtchnPending on the shared
		// info struct.
		if sel&(1<<i) == 0 {
			continue
		}

		basePort := i * 64
		pending := _atman_shared_info.EvtchnPending[i]

		for j := uint32(0); j < 64; j++ {
			if pending&(1<<j) == 0 {
				continue
			}

			handleEvent(basePort+j, r)
		}
	}
}

// handleEvent calls the registered handler on port, if it exists,
// and then acknowledges the event.
func handleEvent(port uint32, r *cpuRegisters) {
	if handler := eventHandlers[port]; handler != nil {
		handler(port, r)
	}

	clearEventChan(port)
}

// eventCallback is called by eventCallbackASM after switching
// to the gsignal stack.
//
//go:nosplit
//go:nowritebarrierrec
func eventCallback(r *cpuRegisters, sp uintptr) {
	g := getg()
	if g == nil {
		kprintString("Event callback called with no g\n")
		crash()
		return
	}

	if sp < g.m.gsignal.stack.lo || sp >= g.m.gsignal.stack.hi {
		kprintString("eventgo: Event callback called, but not on signal stack\n")
		crash()
		return
	}

	setg(g.m.gsignal)

	// Correct CS and SS for return
	r.cs |= 3
	r.ss |= 3

	handleEvents(r)

	if taskrunqueue.Head != nil && taskrunqueue.Head != taskcurrent {
		// there's another runnable task, let's context switch
		taskcurrent.Context.r = *r
		taskready(taskcurrent)
		atomic.Storep1(unsafe.Pointer(&taskcurrent), unsafe.Pointer(taskrunqueue.Head))
		taskrunqueue.Remove(taskcurrent)

		*r = taskcurrent.Context.r
	}

	setg(g)
	atmansettls(taskcurrent.Context.tls)
}

//go:nosplit
func atmansettls(tls uintptr)

// eventCallbackASM is called by Xen when there are events to process.
//
//go:nosplit
func eventCallbackASM()

// eventFailsafe is called by Xen if the event callback fails.
//
//go:nosplit
func eventFailsafe()
