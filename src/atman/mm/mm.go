package mm

import (
	"unsafe"
)

// Page represents a page of allocated memory.
type Page struct {
	// Frame is an opaque identifier of the page's machine frame.
	Frame uintptr
	Size  int

	// Data is a pointer to a page of memory of Size bytes.
	Data unsafe.Pointer
}

// AllocPage allocates a physical page, maps it into memory, and returns
// a Page structure describing the page.
func AllocPage() *Page {
	frame, size, data := runtime_allocPage()

	return &Page{
		Frame: frame,
		Size:  size,
		Data:  data,
	}
}

// runtime_allocPage is provided by the runtime package.
func runtime_allocPage() (frame uintptr, size int, data unsafe.Pointer)

// MapFrames maps the provided machine frames into
// contiguous memory and returns a pointer to the block.
func MapFrames(frames []uintptr) unsafe.Pointer
