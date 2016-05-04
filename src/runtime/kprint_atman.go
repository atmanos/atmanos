package runtime

import "unsafe"

func kprintString(s string) {
	_atman_console.write(bytes(s))
}

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

func kprintInt(v int64) {
	if v < 0 {
		kprintString("-")
		v = -v
	}

	kprintUint(uint64(v))
}

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

func kprintPointer(p unsafe.Pointer) {
	kprintHex(uint64(uintptr(p)))
}
