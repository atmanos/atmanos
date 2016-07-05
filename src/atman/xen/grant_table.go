package xen

import (
	"atman/mm"
	"atman/xen/hypercall"
	"unsafe"
)

const (
	DOMID_SELF = 0x7FF0

	grantTableMaxFrames = 4
	grantTableNrEntries = grantTableMaxFrames * int(0x1000/unsafe.Sizeof(GrantTableEntry{}))

	grantTableNrReserved = 8
)

var grantTable *GrantTable

type GrantTable struct {
	entries []GrantTableEntry

	grefs grefFreeList
}

// MapGrantTable maps Xen's shared grant table into memory
// and returns a pointer to it.
func MapGrantTable() *GrantTable {
	if grantTable == nil {
		grantTable = setupGrantTable(grantTableMaxFrames)
	}

	return grantTable
}

func setupGrantTable(nrFrames int) *GrantTable {
	grantFrames := make([]uintptr, nrFrames)

	loadGrantFrames(grantFrames)

	tablePtr := mm.MapFrames(grantFrames)
	table := (*[grantTableNrEntries]GrantTableEntry)(tablePtr)

	grefs := make(grefFreeList, 0, grantTableNrEntries)

	for i := grantTableNrReserved; i < grantTableNrEntries; i++ {
		grefs.put(Gref(i))
	}

	return &GrantTable{
		entries: table[:],
		grefs:   grefs,
	}
}

func (t *GrantTable) GrantAccess(domid uint16, frame uintptr, readOnly bool) (Gref, bool) {
	gref, ok := t.grefs.get()
	if !ok {
		return gref, ok
	}

	flags := hypercall.GTF_permit_access
	if readOnly {
		flags |= hypercall.GTF_readonly
	}

	t.entries[gref].Frame = uint32(frame)
	t.entries[gref].DomID = domid
	MemoryBarrierWrite()
	t.entries[gref].Flags = uint16(flags)

	return gref, true
}

func (t *GrantTable) EndAccess(gref Gref) bool {
	for {
		flags := t.entries[gref].Flags

		if flags&(hypercall.GTF_reading|hypercall.GTF_writing) != 0 {
			return false
		}

		// TODO: compareAndSwapUint16(&t.entries[gref].Flags, flags, 0)
		t.entries[gref].Flags = 0
		break
	}

	t.grefs.put(gref)
	return true
}

func loadGrantFrames(frames []uintptr) uintptr {
	setup := struct {
		DomID    uint16
		_        uint16
		NrFrames uint32
		Status   uint16
		_        [6]byte
		Frames   *uintptr
	}{
		DomID:    DOMID_SELF,
		NrFrames: uint32(len(frames)),
		Frames:   &frames[0],
	}

	return hypercall.GrantTableOp(
		hypercall.GNTTABOP_setup_table,
		unsafe.Pointer(&setup),
	)
}

// GrantTableEntry is the V2 grant table entry structure for full pages.
type GrantTableEntry struct {
	Flags uint16
	DomID uint16
	Frame uint32
}

type Gref uint32

// grefFreeList is a free-list of grefs, where each gref
// is an index to an unused entry in the grant table.
type grefFreeList []Gref

func (l *grefFreeList) get() (Gref, bool) {
	s := *l
	if len(s) == 0 {
		return 0, false
	}
	gref := s[len(s)-1]
	*l = s[:len(s)-1]
	return gref, true
}

func (l *grefFreeList) put(gref Gref) {
	old := *l
	*l = append(old, gref)
}
