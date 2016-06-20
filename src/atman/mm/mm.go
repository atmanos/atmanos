package mm

import (
	"unsafe"
)

// Page represents a page of allocated memory.
type Page struct {
	// Frame is an opaque identifier of the page's machine frame.
	Frame uintptr
	// Ptr is a pointer to a page of memory of size len(Data).
	Ptr unsafe.Pointer
	// Data is the page represented as a slice of bytes.
	Data []byte
}

// AllocPage allocates a physical page, maps it into memory, and returns
// a Page structure describing the page.
func AllocPage() *Page {
	frame, data := runtime_allocPage()

	return &Page{
		Frame: frame,
		Data:  data,
		Ptr:   unsafe.Pointer(&data[0]),
	}
}

// runtime_allocPage is provided by the runtime package.
func runtime_allocPage() (frame uintptr, data []byte)

// MapFrames maps the provided machine frames into
// contiguous memory and returns a pointer to the block.
func MapFrames(frames []uintptr) unsafe.Pointer
