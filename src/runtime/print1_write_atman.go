package runtime

import "unsafe"

func writeErr(b []byte) {
	HYPERVISOR_console_io(
		0,
		uint64(len(b)),
		uintptr(unsafe.Pointer(&b[0])),
	)
}
