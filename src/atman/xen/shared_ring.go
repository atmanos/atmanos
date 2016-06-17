package xen

import (
	"unsafe"
)

// SharedRing is a shared-memory producer-consumer ring.
//
// It's memory layout is precise, and should be accessed
// through FrontendRing.
type SharedRing struct {
	RequestProducer uint32
	RequestEvent    uint32

	ResponseProducer uint32
	ResponseEvent    uint32

	_ uint64    // private union data for backend
	_ [44]uint8 // pad to 512 bytes

	// The remainder of the ring page is available for data.
	Data [3584]byte
}

// entryCount returns the number of entries the ring can store based on
// entrySize.
func (r *SharedRing) entryCount(entrySize int) int {
	return r.round(len(r.Data) / entrySize)
}

// round rounds c down to the nearest power of 2.
func (r *SharedRing) round(x int) int {
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return (x >> 1) + 1
}

// FrontendRing provides the frontend interface to a SharedRing.
type FrontendRing struct {
	RequestProducer  uint32
	ResponseConsumer uint32

	EntrySize  uint32
	EntryCount uint32

	*SharedRing
}

// NewFrontendRing initializes a new FrontendRing from the provided SharedRing.
//
// entrySize should be the size in bytes of the largest value to be stored in
// the ring.
func NewFrontendRing(r *SharedRing, entrySize int) *FrontendRing {
	return &FrontendRing{
		EntrySize:  uint32(entrySize),
		EntryCount: uint32(r.entryCount(entrySize)),
		SharedRing: r,
	}
}

// NextRequest returns a pointer to the next request,
// advancing the private producer index.
func (r *FrontendRing) NextRequest() unsafe.Pointer {
	req := r.get(r.RequestProducer)
	r.RequestProducer++

	return req
}

// get returns a pointer to the entry at index i.
func (r *FrontendRing) get(i uint32) unsafe.Pointer {
	offset := i & (r.EntryCount - 1)
	idx := offset * r.EntrySize

	return unsafe.Pointer(&r.SharedRing.Data[idx])
}

// PushRequests updates the shared ring so the backend
// sees all pending requests.
//
// If notify is true, the backend expects to receive
// an event channel notification.
func (r *FrontendRing) PushRequests() (notify bool) {
	old := r.SharedRing.RequestProducer
	new := r.RequestProducer

	MemoryBarrierWrite() // backend sees requests before we update index.

	r.SharedRing.RequestProducer = new

	MemoryBarrier() // backend sees new requests before check notify

	event := r.SharedRing.RequestEvent

	// If RequestEvent is ahead of old, the backend consumed all
	// pending requests and has requested a notification for new events.
	//
	// The expression is equivalent to event > old,
	// but handles overflow of the index type.
	return (new - event) < (new - old)
}
