package ring

import "sync/atomic"

// Buffer provides a memory-layout agnostic and safe interface
// for writing to or reading from a shared Xen ring buffer.
type Buffer struct {
	Data []byte

	ReaderPos *uint32
	WriterPos *uint32
}

// Write writes bytes from data info buf, updates the position,
// and returns the number of bytes written.
func (buf *Buffer) Write(data []byte) int {
	var (
		bufSize = uint32(len(buf.Data))
		mask    = bufSize - 1

		reader = atomic.LoadUint32(buf.ReaderPos)
		writer = atomic.LoadUint32(buf.WriterPos)

		avail = int(bufSize - (writer - reader))
		size  = len(data)
	)

	if size > avail {
		size = avail
	}

	for i := 0; i < size; i++ {
		buf.Data[writer&mask] = data[i]
		writer++
	}

	atomic.StoreUint32(buf.WriterPos, writer)
	return size
}

// Read reads bytes from buf into data, updates the position,
// and returns the number of bytes read.
func (buf *Buffer) Read(data []byte) int {
	var (
		bufSize = uint32(len(buf.Data))
		mask    = bufSize - 1

		reader = atomic.LoadUint32(buf.ReaderPos)
		writer = atomic.LoadUint32(buf.WriterPos)

		size = int(writer - reader)
	)

	if size > len(data) {
		size = len(data)
	}

	for i := 0; i < size; i++ {
		data[i] = buf.Data[reader&mask]
		reader++
	}

	atomic.StoreUint32(buf.ReaderPos, reader)
	return size
}
