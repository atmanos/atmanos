package runtime

import "unsafe"

var _atman_console console

type console struct {
	port uint32

	*consoleRing
}

func (c *console) init() {
	c.port = _atman_start_info.Console.Eventchn
	c.consoleRing = (*consoleRing)(unsafe.Pointer(
		_atman_start_info.Console.Mfn.pfn().vaddr(),
	))
}

func (c console) notify() {
	eventChanSend(c.port)
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
		if consoleRingOutSize-prod-cons == 0 {
			break
		}

		i := prod & (consoleRingOutSize - 1)
		r.out[i] = c

		prod++
		sent++
	}

	atomicstore(&r.outProducerPos, prod)
	return sent
}

//go:linkname syscall_WriteConsole syscall.WriteConsole
func syscall_WriteConsole(b []byte) int {
	n := int(_atman_console.write(b))
	_atman_console.notify()
	return n
}
