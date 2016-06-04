package runtime

import (
	"runtime/internal/atomic"
)

// NotifyEventChannel sends an event signal to the listener on port.
func NotifyEventChannel(port uint32) {
	eventChanSend(port)
}

// BindEventHandler registers intent to receive events on port
// by later calling WaitEvent.
func BindEventHandler(port uint32) {
	noteclear(&eventNotifications[port].note)

	bindEventHandler(port, notifyEvent)

	unmaskEventChan(port)
}

// WaitEvent waits for an event to fire on port.
//
// Only one goroutine should call WaitEvent.
func WaitEvent(port uint32) {
	n := &eventNotifications[port]
	n.recv()
}

func notifyEvent(port uint32, _ *cpuRegisters) {
	n := &eventNotifications[port]
	n.send()
}

// Below is an adaptation of the signal handling mechanisms
// from runtime/sigqueue.go.

var eventNotifications [1024]eventNotification

const (
	eventNotificationIdle = iota
	eventNotificationReceiving
	eventNotificationSending
)

type eventNotification struct {
	note  note
	state uint32
}

func (n *eventNotification) recv() {
	for {
		switch atomic.Load(&n.state) {
		default:
			throw("eventNotification.recv: inconsistent state")
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

func (n *eventNotification) send() {
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
