package runtime

import "unsafe"

// clearbit atomically clears the bit at offset n from addr.
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

func unmaskEventChan(port uint32) {
	info := _atman_shared_info
	clearbit(&info.EvtchnMask[0], port)
}
