package runtime

import (
	"unsafe"
)

// The memory management functions required by the runtime.

func sysFree(v unsafe.Pointer, n uintptr, sysStat *uint64) {}

func sysUnused(v unsafe.Pointer, n uintptr) {}

func sysUsed(v unsafe.Pointer, n uintptr) {}

func sysFault(v unsafe.Pointer, n uintptr) {}

// sysReserve reserves n bytes at v and updates reserved.
func sysReserve(v unsafe.Pointer, n uintptr, reserved *bool) unsafe.Pointer {
	*reserved = false
	return v
}

// sysMap makes n bytes at v readable and writable and adjusts the stats.
//go:nosplit
func sysMap(v unsafe.Pointer, n uintptr, reserved bool, sysStat *uint64) {
	mSysStatInc(sysStat, n)
	p := _atman_mm.allocN(v, n)
	if p == nil {
		kprintString("runtime: out of memory\n")
		crash()
	}
	if p != v {
		throw("runtime: cannot map pages in arena address space")
	}
}

// sysAlloc allocates n bytes, adjusts sysStat, and returns the address
// of the allocated bytes.
//go:nosplit
func sysAlloc(n uintptr, sysStat *uint64) unsafe.Pointer {
	p := _atman_mm.allocN(nil, n)
	if p != nil {
		mSysStatInc(sysStat, n)
	}
	return p
}
