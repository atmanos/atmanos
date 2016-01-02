package runtime

import "unsafe"

func sysFree(v unsafe.Pointer, n uintptr, sysStat *uint64) {}

func sysUnused(v unsafe.Pointer, n uintptr) {}
func sysUsed(v unsafe.Pointer, n uintptr)   {}
func sysFault(v unsafe.Pointer, n uintptr)  {}

// sysMap makes n bytes at v readable and writable and adjusts the stats.
func sysMap(v unsafe.Pointer, n uintptr, reserved bool, sysStat *uint64) {
	mSysStatInc(sysStat, n)
	p := memAlloc(v, n)
	if p != v {
		throw("runtime: cannot map pages in arena address space")
	}
}

// sysAlloc allocates n bytes, adjusts sysStat, and returns the address
// of the allocated bytes.
func sysAlloc(n uintptr, sysStat *uint64) unsafe.Pointer {
	p := memAlloc(nil, n)
	if p != nil {
		mSysStatInc(sysStat, n)
	}
	return p
}

// sysReserve reserves n bytes at v and updates reserved.
func sysReserve(v unsafe.Pointer, n uintptr, reserved *bool) unsafe.Pointer {
	*reserved = false
	return v
}

// memAlloc allocates n bytes of memory at address v
// and returns a pointer to the allocated memory.
// If v is nil, an address will be chosen.
func memAlloc(v unsafe.Pointer, n uintptr) unsafe.Pointer {
	requiredPages := uint64(round(n, _PAGESIZE) / _PAGESIZE)

	return _atman_mm.allocPages(v, requiredPages)
}

var _atman_mm = &atmanMemoryManager{}

type atmanMemoryManager struct {
	bootstrapStackPFN pfn // start of bootstrap stack
	bootstrapEndPFN   pfn // end of bootstrap region

	nextPFN pfn // next free frame
	lastPFN pfn

	nextHeapPage vaddr

	l4PFN pfn
	l4    xenPageTable
}

func (mm *atmanMemoryManager) init() {
	var (
		pageTableBase = _atman_start_info.PageTableBase
		ptStartPfn    = pageTableBase.pfn()
		ptEndPfn      = ptStartPfn.add(_atman_start_info.NrPageTableFrames)

		bootstrapStackPFN  = ptEndPfn.add(1)
		bootstrapStackAddr = bootstrapStackPFN.vaddr()

		bootstrapEnd = round(
			uintptr(bootstrapStackAddr)+0x80000, // minimum 512kB padding
			0x400000, // 4MB alignment
		)
		bootstrapEndPFN = vaddr(bootstrapEnd).pfn()
	)

	mm.bootstrapStackPFN = bootstrapStackPFN
	mm.bootstrapEndPFN = bootstrapEndPFN
	mm.nextPFN = bootstrapEndPFN.add(1)
	mm.lastPFN = pfn(_atman_start_info.NrPages)

	mm.nextHeapPage = mm.nextPFN.vaddr()

	mm.l4PFN = pageTableBase.pfn()
	mm.l4 = mm.mapL4(mm.l4PFN)

	mm.unmapLowAddresses()
	mm.unmapBootstrapPageTables()
}

func (mm *atmanMemoryManager) unmapLowAddresses() {
	addr := 0

	for addr < 0x40000 {
		HYPERVISOR_update_va_mapping(uintptr(addr), 0, 2)
		addr += _PAGESIZE
	}
}

func (mm *atmanMemoryManager) allocPages(v unsafe.Pointer, n uint64) unsafe.Pointer {
	if v == nil {
		v = mm.reserveHeapPages(n)
	}

	for page := vaddr(v); page < vaddr(v)+vaddr(n*_PAGESIZE); page += _PAGESIZE {
		mm.allocPage(page)
	}

	return v
}

// allocPage makes page a writeable userspace page.
func (mm *atmanMemoryManager) allocPage(page vaddr) {
	var (
		l4offset = page.pageTableOffset(pageTableLevel4)
		l3offset = page.pageTableOffset(pageTableLevel3)
		l2offset = page.pageTableOffset(pageTableLevel2)
		l1offset = page.pageTableOffset(pageTableLevel1)
	)

	l4 := mm.l4
	l3pte := l4.Get(l4offset)

	if !l3pte.hasFlag(xenPageTablePresent) {
		pfn := mm.physAllocPage()
		l3pte = mm.writePte(mm.l4PFN, l4offset, pfn, PTE_PAGE_TABLE_FLAGS|xenPageTableWritable)
	}

	l3 := mm.getPageTable(-1, -1, l4offset)
	l2pte := l3.Get(l3offset)

	if !l2pte.hasFlag(xenPageTablePresent) {
		pfn := mm.physAllocPage()
		l2pte = mm.writePte(l3pte.pfn(), l3offset, pfn, PTE_PAGE_TABLE_FLAGS|xenPageTableWritable)
	}

	l2 := mm.getPageTable(-1, l4offset, l3offset)
	l1pte := l2.Get(l2offset)

	if !l1pte.hasFlag(xenPageTablePresent) {
		pfn := mm.physAllocPage()
		l1pte = mm.writePte(l2pte.pfn(), l2offset, pfn, PTE_PAGE_TABLE_FLAGS|xenPageTableWritable)
	}

	pagepfn := mm.physAllocPage()
	mm.writePte(l1pte.pfn(), l1offset, pagepfn, PTE_PAGE_FLAGS)

	// ensure page is writable
	*(*uintptr)(unsafe.Pointer(page)) = 0x0
}

func (mm *atmanMemoryManager) getPageTable(a, b, c int) xenPageTable {
	const pageTableVaddrOffset = vaddr(0xFFFFFF8000000000)

	if a == -1 {
		a = 511
	}

	if b == -1 {
		b = 511
	}

	if c == -1 {
		c = 511
	}

	addr := pageTableVaddrOffset + vaddr(a<<30) + vaddr(b<<21) + vaddr(c<<12)

	return newXenPageTable(addr)
}

func (mm *atmanMemoryManager) physAllocPage() pfn {
	pfn := mm.reservePFN()
	mm.clearPage(pfn)
	return pfn
}

func (mm *atmanMemoryManager) pageTableWalk(addr vaddr) {
	var (
		l4offset = addr.pageTableOffset(pageTableLevel4)
		l3offset = addr.pageTableOffset(pageTableLevel3)
		l2offset = addr.pageTableOffset(pageTableLevel2)
		l1offset = addr.pageTableOffset(pageTableLevel1)

		l4 = mm.l4
	)

	println("page table walk from", unsafe.Pointer(addr))
	print("L4[")
	print(l4offset)
	print("] = ")

	l3pte := l4.Get(l4offset)
	l3pte.debug()

	if !l3pte.hasFlag(xenPageTablePresent) {
		return
	}

	l3 := mm.getPageTable(-1, -1, l4offset)
	print("L3[")
	print(l3offset)
	print("] = ")

	l2pte := l3.Get(l3offset)
	l2pte.debug()

	if !l2pte.hasFlag(xenPageTablePresent) {
		return
	}

	l2 := mm.getPageTable(-1, l4offset, l3offset)
	print("L2[")
	print(l2offset)
	print("] = ")

	l1pte := l2.Get(l2offset)
	l1pte.debug()

	if !l1pte.hasFlag(xenPageTablePresent) {
		return
	}

	l1 := mm.getPageTable(l4offset, l3offset, l2offset)
	print("L1[")
	print(l1offset)
	print("] = ")

	l0pte := l1.Get(l1offset)
	l0pte.debug()

	if !l0pte.hasFlag(xenPageTablePresent) {
		return
	}
}

func (mm *atmanMemoryManager) reserveHeapPages(n uint64) unsafe.Pointer {
	var p vaddr
	p, mm.nextHeapPage = mm.nextHeapPage, mm.nextHeapPage+vaddr(n*_PAGESIZE)
	return unsafe.Pointer(p)
}

func (mm *atmanMemoryManager) reservePFN() pfn {
	var p pfn
	p, mm.nextPFN = mm.nextPFN, mm.nextPFN+1
	return p
}

// mapL4 sets up recursively mapped page table
// from the initial bootstrap page tables.
func (mm *atmanMemoryManager) mapL4(pfn pfn) xenPageTable {
	mm.writePte(pfn, 511, pfn, PTE_PAGE_TABLE_FLAGS)

	return mm.getPageTable(-1, -1, -1)
}

// unmapBootstrapPageTables removes the bootstrap page table mappings.
func (mm *atmanMemoryManager) unmapBootstrapPageTables() {
	for i := uint64(0); i < _atman_start_info.NrPageTableFrames; i++ {
		addr := mm.l4PFN.add(i).vaddr()
		HYPERVISOR_update_va_mapping(uintptr(addr), 0, 2)
	}
}

func (mm *atmanMemoryManager) clearPage(pfn pfn) {
	mm.mmuExtOp([]mmuExtOp{
		{
			cmd:  16, // MMUEXT_CLEAR_PAGE
			arg1: uint64(pfn.mfn()),
		},
	})
}

func (mm *atmanMemoryManager) mmuExtOp(ops []mmuExtOp) {
	ret := HYPERVISOR_mmuext_op(ops, DOMID_SELF)

	if ret != 0 {
		println("HYPERVISOR_mmuext_op returned", ret)
	}
}

func (mm *atmanMemoryManager) writePte(table pfn, offset int, value pfn, flags uintptr) pageTableEntry {
	newpte := pageTableEntry(value.mfn() << xenPageFlagShift)
	newpte.setFlag(flags)

	updates := []mmuUpdate{
		{
			ptr: uintptr((table.mfn() << xenPageFlagShift)) + uintptr(offset*ptrSize),
			val: uintptr(newpte),
		},
	}
	ret := HYPERVISOR_mmu_update(updates, DOMID_SELF)

	if ret != 0 {
		println("writePte: HYPERVISOR_mmu_update returned", ret)
	}

	return newpte
}
