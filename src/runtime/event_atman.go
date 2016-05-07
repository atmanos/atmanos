package runtime

import "unsafe"

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

//go:nosplit
func hypervisorEventCallback()

//go:nosplit
func hypervisorFailsafeCallback()

//go:nosplit
func handleHypervisorCallback(r *cpuRegisters) {
	_atman_shared_info.VCPUInfo[0].UpcallPending = 0
	_atman_shared_info.VCPUInfo[0].PendingSel = 0
	_atman_shared_info.EvtchnPending[0] = 0
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
