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

const (
	EINVAL = Errno(iota + 1)
	EISDIR
	ENOTDIR
	ENAMETOOLONG
	EPIPE
)

const (
	O_RDONLY = 1 << iota
	O_WRONLY
	O_RDWR
	O_CREAT
	O_APPEND
	O_TRUNC
	O_EXCL
	O_SYNC
)

var (
	Stdin  = 0
	Stdout = 1
	Stderr = 2
)

type Timespec struct {
	Sec  int64
	Nsec int64
}

type Timeval struct {
	Sec  int64
	Usec int64
}

type SysProcAttr struct{}

func Getenv(s string) (string, bool) {
	return "", false
}

func Setenv(key, value string) error {
	return nil
}

func Unsetenv(key string) error {
	return nil
}

func Environ() []string {
	return nil
}

func Clearenv() {}

func Getpagesize() int { return 0x1000 }

func Getppid() int { return 2 }

func Getpid() int { return 3 }

func Getuid() int               { return 1 }
func Geteuid() int              { return 1 }
func Getgid() int               { return 1 }
func Getegid() int              { return 1 }
func Getgroups() ([]int, error) { return []int{1}, nil }

const ImplementsGetwd = false

func Getwd() (dir string, err error) {
	return "", EINVAL
}

func Chdir(path string) error {
	return EINVAL
}

func Fchdir(fd int) error {
	return EINVAL
}

func Mkdir(path string, perm uint32) error {
	return EINVAL
}

func Exit(code int) {}

// Write writes the contents of b to fd and returns
// the number of bytes written or an error.
//
// If fd is Stdout, b is written with WriteConsole.
func Write(fd int, b []byte) (n int, err error) {
	if fd != 1 {
		return 0, EINVAL
	}

	return WriteConsole(b), nil
}

// Read is currently unsupported.
func Read(fd int, b []byte) (n int, err error) {
	return 0, EINVAL
}

// WriteConsole writes b to the Xen console and
// returns the number of bytes written.
func WriteConsole(b []byte) int
