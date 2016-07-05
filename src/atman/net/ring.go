package net

import (
	"atman/mm"
	"atman/xen"
	"unsafe"
)

type netifTxRequest struct {
	Gref   xen.Gref
	Offset uint16
	Flags  uint16
	ID     uint16
	Size   uint16
}

type netifTxResponse struct {
	ID     uint16
	Status int16
}

func newTxRing(page *mm.Page) *xen.FrontendRing {
	return xen.NewFrontendRing(
		(*xen.SharedRing)(page.Ptr),
		txRingEntrySize(),
	)
}

func txRingEntrySize() int {
	size := unsafe.Sizeof(netifTxRequest{})

	if rspSize := unsafe.Sizeof(netifTxResponse{}); rspSize > size {
		size = rspSize
	}

	return int(size)
}

type netifRxRequest struct {
	ID   uint16
	_    uint16
	Gref xen.Gref
}

type netifRxResponse struct {
	ID     uint16
	Offset uint16
	Flags  uint16
	Status int16
}

func newRxRing(page *mm.Page) *xen.FrontendRing {
	return xen.NewFrontendRing(
		(*xen.SharedRing)(page.Ptr),
		rxRingEntrySize(),
	)
}

func rxRingEntrySize() int {
	size := unsafe.Sizeof(netifRxRequest{})

	if rspSize := unsafe.Sizeof(netifRxResponse{}); rspSize > size {
		size = rspSize
	}

	return int(size)
}
