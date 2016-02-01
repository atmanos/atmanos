package os

import "syscall"

type File struct {
	fd      int
	name    string
	dirinfo *dirInfo
}

type dirInfo struct{}

func NewFile(fd uintptr, name string) *File {
	return &File{
		fd:   int(fd),
		name: name,
	}
}

func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	return nil, ErrNotExist
}

func (*File) Close() error {
	return ErrInvalid
}

func (*File) Stat() (FileInfo, error) {
	return nil, ErrInvalid
}

func (*File) read(b []byte) (n int, err error) {
	return 0, ErrNotExist
}

func (*File) pread(b []byte, off int64) (n int, err error) {
	return 0, ErrNotExist
}

func (f *File) write(b []byte) (n int, err error) {
	return syscall.Write(f.fd, b)
}

func (*File) pwrite(b []byte, off int64) (n int, err error) {
	return 0, ErrNotExist
}

func (*File) readdir(n int) (fi []FileInfo, err error) {
	return nil, ErrNotExist
}

func (*File) readdirnames(n int) (names []string, err error) {
	return nil, ErrNotExist
}

func (*File) seek(offset int64, whence int) (ret int64, err error) {
	return 0, ErrNotExist
}

func epipecheck(file *File, e error) {
}

func syscallMode(i FileMode) (o uint32) {
	return 0
}

func Stat(name string) (FileInfo, error) {
	return nil, ErrNotExist
}

func Lstat(name string) (FileInfo, error) {
	return nil, ErrNotExist
}

func Chmod(name string, mode FileMode) error {
	return ErrNotExist
}

func Remove(path string) error {
	return ErrNotExist
}

func rename(oldname, newname string) error {
	return ErrNotExist
}

func sameFile(fs1, fs2 *fileStat) bool {
	return false
}
