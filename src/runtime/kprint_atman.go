package runtime

import "unsafe"

//go:nowritebarrier
func kprintString(s string) {
	var buf [100]byte
	copy(buf[:], s)
	_atman_console.write(buf[:len(s)])
}

//go:nowritebarrier
func kprintUint(v uint64) {
	var buf [100]byte
	i := len(buf)
	for i--; i > 0; i-- {
		buf[i] = byte(v%10 + '0')
		if v < 10 {
			break
		}
		v /= 10
	}
	_atman_console.write(buf[i:])
}

//go:nowritebarrier
func kprintInt(v int64) {
	if v < 0 {
		kprintString("-")
		v = -v
	}

	kprintUint(uint64(v))
}

//go:nowritebarrier
func kprintHex(v uint64) {
	const dig = "0123456789abcdef"
	var buf [100]byte
	i := len(buf)
	for i--; i > 0; i-- {
		buf[i] = dig[v%16]
		if v < 16 {
			break
		}
		v /= 16
	}
	i--
	buf[i] = 'x'
	i--
	buf[i] = '0'
	_atman_console.write(buf[i:])
}

//go:nowritebarrier
func kprintPointer(p unsafe.Pointer) {
	kprintHex(uint64(uintptr(p)))
}
