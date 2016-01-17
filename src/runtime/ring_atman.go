package runtime

import "unsafe"

type netifTxResponse struct {
	id     uint16
	status int16
}

type netifTxRequest struct {
	gref   uint32 // grantRef
	offset uint16
	flags  uint16
	id     uint16
	size   uint16
}

type ringIndex uint64

// netifTxSharedRingEntry can be one of
// netifTxRequest or netifTxResponse.
type netifTxSharedRingEntry [16]uint8

func (e *netifTxSharedRingEntry) req() *netifTxRequest {
	return (*netifTxRequest)(unsafe.Pointer(e))
}

const (
	netifTxSharedRingEntrySize   = 128
	netifTxSharedRingEntryOffset = 512
	netifTxSharedRingNrEntries   = (_PAGESIZE - netifTxSharedRingEntryOffset) / netifTxSharedRingEntrySize
)

type netifTxSharedRing struct {
	requestProducer  ringIndex // unsigned int
	requestEvent     ringIndex
	responseProducer ringIndex
	responseEvent    ringIndex
	_                uint64    // private union data
	_                [44]uint8 // pad to 512 bytes
	ring             [netifTxSharedRingNrEntries]netifTxSharedRingEntry
}

type netifTxFrontRing struct {
	requestProducer  ringIndex
	responseConsumer ringIndex
	nrEntries        uint
	shared           *netifTxSharedRing
}

func (r *netifTxFrontRing) getRequest() *netifTxRequest {
	// return r.shared.ring.get(r.requestProducer).req()
	return r.shared.ring[r.requestProducer&(netifTxSharedRingNrEntries-1)].req()
}

func (r *netifTxFrontRing) pushRequests() {
	atomicstore64(
		(*uint64)(unsafe.Pointer(&r.shared.requestProducer)),
		uint64(r.requestProducer),
	)
}
