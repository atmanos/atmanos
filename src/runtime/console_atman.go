package runtime

import "unsafe"

var _atman_console console

type console struct {
	port uint32

	ring *consoleRing
}

func (c *console) init() {
	c.port = _atman_start_info.Console.Eventchn
	c.ring = (*consoleRing)(unsafe.Pointer(
		_atman_start_info.Console.Mfn.pfn().vaddr(),
	))
}

func (c console) notify() {
	eventChanSend(c.port)
}

func (c console) write(b []byte) int {
	n := c.ring.write(b)
	c.notify()
	return int(n)
}

const (
	consoleRingInSize  = 1024
	consoleRingOutSize = 2048
)

type consoleRing struct {
	in  [consoleRingInSize]byte
	out [consoleRingOutSize]byte

	inConsumerPos uint32
	inProducerPos uint32

	outConsumerPos uint32
	outProducerPos uint32
}

func (r *consoleRing) write(b []byte) uint32 {
	var (
		sent = uint32(0)

		cons = atomicload(&r.outConsumerPos)
		prod = atomicload(&r.outProducerPos)
	)

	for _, c := range b {
		size := uint32(1)
		if c == '\n' {
			size = 2
		}

		if consoleRingOutSize-prod-cons < size {
			break
		}

		if c == '\n' {
			r.writeByteAt('\r', prod)
			prod++
			sent++
		}

		r.writeByteAt(c, prod)
		prod++
		sent++
	}

	atomicstore(&r.outProducerPos, prod)
	return sent
}

func (r *consoleRing) writeByteAt(b byte, off uint32) {
	i := off & (consoleRingOutSize - 1)
	r.out[i] = b
}

//go:linkname syscall_WriteConsole syscall.WriteConsole
func syscall_WriteConsole(b []byte) int {
	return _atman_console.write(b)
}
