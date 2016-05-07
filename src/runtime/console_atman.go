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
	bindEventHandler(c.port, consoleHandleInput)
}

// consoleHandleInput is run when console input is available.
// It reads from the console, and echoes the data back.
func consoleHandleInput(_ uint32, _ *cpuRegisters) {
	var buf [100]byte

	for {
		n := _atman_console.read(buf[:])
		if n == 0 {
			break
		}

		_atman_console.write(buf[:n])

		if buf[n-1] == '\r' {
			_atman_console.write([]byte{'\n'})
		}
	}
}

func (c console) notify() {
	eventChanSend(c.port)
}

func (c console) write(b []byte) int {
	n := c.ring.write(b)
	c.notify()
	return int(n)
}

func (c console) read(b []byte) int {
	n := c.ring.read(b)
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

func (r *consoleRing) read(b []byte) int {
	var (
		cons = atomicload(&r.inConsumerPos)
		prod = atomicload(&r.inProducerPos)
	)

	size := int(prod) - int(cons)
	if size > len(b) {
		size = len(b)
	}

	for i := 0; i < size; i++ {
		b[i] = r.in[cons&(consoleRingInSize-1)]
		cons++
	}

	atomicstore(&r.inConsumerPos, cons)
	return size
}

//go:linkname syscall_WriteConsole syscall.WriteConsole
func syscall_WriteConsole(b []byte) int {
	return _atman_console.write(b)
}
