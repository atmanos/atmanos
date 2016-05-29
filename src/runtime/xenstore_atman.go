package runtime

import (
	"runtime/internal/atomic"
	"unsafe"
)

func LoadXenStore() (ring unsafe.Pointer, port uint32) {
	port = _atman_start_info.Store.Eventchn
	ring = unsafe.Pointer(_atman_start_info.Store.Mfn.pfn().vaddr())

	return ring, port
}

var eventNotifications [1024]struct {
	note  note
	state uint32
}

const (
	eventNotificationIdle = iota
	eventNotificationReceiving
	eventNotificationSending
)

func WaitEvent(port uint32) {
	n := &eventNotifications[port]

	for {
		switch atomic.Load(&n.state) {
		default:
			throw("eventNotificationnal_recv: inconsistent state")
		case eventNotificationIdle:
			if atomic.Cas(&n.state, eventNotificationIdle, eventNotificationReceiving) {
				notetsleepg(&n.note, -1)
				noteclear(&n.note)
				return
			}
		case eventNotificationSending:
			if atomic.Cas(&n.state, eventNotificationSending, eventNotificationIdle) {
				return
			}
		}
	}
}

func NotifyEventChannel(port uint32) {
	eventChanSend(port)
}

func BindEventHandler(port uint32) {
	noteclear(&eventNotifications[port].note)

	bindEventHandler(port, notifyEvent)

	unmaskEventChan(port)
}

func notifyEvent(port uint32, _ *cpuRegisters) {
	n := &eventNotifications[port]

	for {
		switch atomic.Load(&n.state) {
		default:
			throw("eventNotificationsend: inconsistent state")
		case eventNotificationIdle:
			if atomic.Cas(&n.state, eventNotificationIdle, eventNotificationSending) {
				return
			}
		case eventNotificationSending:
			// notification already pending
			return
		case eventNotificationReceiving:
			if atomic.Cas(&n.state, eventNotificationReceiving, eventNotificationIdle) {
				notewakeup(&n.note)
				return
			}
		}
	}
}
