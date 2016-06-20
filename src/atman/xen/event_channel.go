package xen

import (
	"atman/xen/hypercall"
	"runtime"
	"unsafe"
)

type EventChannel struct {
	Port uint32
}

func NewEventChannel(remote uint32) *EventChannel {
	ec := &EventChannel{}

	ec.allocPort(remote)
	runtime.BindEventHandler(ec.Port)

	return ec
}

func (ec *EventChannel) Notify() {
	runtime.NotifyEventChannel(ec.Port)
}

func (ec *EventChannel) Wait() {
	runtime.WaitEvent(ec.Port)
}

func (ec *EventChannel) allocPort(remote uint32) {
	op := struct {
		domID       uint16
		remoteDomID uint16
		port        uint32
	}{
		domID:       DOMID_SELF,
		remoteDomID: uint16(remote),
	}

	ret := hypercall.EventChannelOp(
		hypercall.EVTCHNOP_alloc_unbound,
		unsafe.Pointer(&op),
	)

	if ret != 0 {
		println("allocPort: alloc event channel failed with ret=", ret)
		panic("allocPort: unable to alloc event channel")
	}

	ec.Port = op.port
}
