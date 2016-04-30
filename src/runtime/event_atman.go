package runtime

import "unsafe"

type eventHandler func(port uint32, r *cpuRegisters)

var eventHandlers = make([]eventHandler, 4096)

func initEvents() {
	HYPERVISOR_set_callbacks(
		funcPC(hypervisorEventCallback),
		funcPC(hypervisorFailsafeCallback),
	)

	enableIRQ()
	bindVIRQ()
}

// clearbit atomically clears the bit at offset n from addr.
//
//go:nosplit
func clearbit(addr *uint64, n uint32)

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
	info := _atman_shared_info
	clearbit(&info.EvtchnPending[0], port)
}

// enableIRQ enables interrupts clearing the event mask
// on the shared VCPU structure.
//
//go:nosplit
func enableIRQ() {
	vcpu := &_atman_shared_info.VCPUInfo[0]
	vcpu.UpcallMask = 0
}

//go:nosplit
func hypervisorEventCallback()

//go:nosplit
func hypervisorFailsafeCallback()

//go:nosplit
func handleHypervisorCallback(r *cpuRegisters) {
	_atman_shared_info.VCPUInfo[0].UpcallPending = 0

	systemstack(func() {
		sel := xchg(&_atman_shared_info.VCPUInfo[0].PendingSel, 0)

		for i := uint32(0); i < 64; i++ {
			// each set bit in sel is an index into EvtchnPending on the shared
			// info struct.
			if sel&(1<<i) == 0 {
				continue
			}

			port := i * 64
			pending := _atman_shared_info.EvtchnPending[sel]

			for j := uint32(0); j < 64; j++ {
				if pending&(1<<j) == 0 {
					continue
				}
				port += j

				handleEvent(port, r)
			}
		}
	})
}

func handleEvent(port uint32, r *cpuRegisters) {
	if handler := eventHandlers[port]; handler != nil {
		handler(port, r)
	}

	clearEventChan(port)
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

	println("Bound VIRQ to Port", op.Port)
}

// bindEventHandler binds the handler f to events on port.
func bindEventHandler(port uint32, f eventHandler) {
	eventHandlers[port] = f
}
