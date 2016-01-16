package runtime

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
