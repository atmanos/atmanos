package mm

import (
	"unsafe"
)

// MapFrames maps the provided machine frames into
// contiguous memory and returns a pointer to the block.
func MapFrames(frames []uintptr) unsafe.Pointer
