package syscall

// An Errno is an unsigned number describing an error condition.
// It implements the error interface.  The zero Errno is by convention
// a non-error, so code to convert from Errno to error should use:
//	err = nil
//	if errno != 0 {
//		err = errno
//	}
type Errno uintptr

func (e Errno) Error() string {
	return "errno " + itoa(int(e))
}

const EINVAL = Errno(0x1)

type Timespec struct {
	Sec  int64
	Nsec int64
}

type Timeval struct {
	Sec  int64
	Usec int64
}

func Getenv(s string) (string, error) {
	return "", EINVAL
}
