package runtime

import (
	"unsafe"
)

// LoadXenStore returns a pointer to the shared-memory
// XenStore ring interface and event channel port.
func LoadXenStore() (ring unsafe.Pointer, port uint32) {
	port = _atman_start_info.Store.Eventchn
	ring = unsafe.Pointer(_atman_start_info.Store.Mfn.pfn().vaddr())

	return ring, port
}
