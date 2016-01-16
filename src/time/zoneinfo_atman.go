package time

import "errors"

var errNotImplemented = errors.New("Not implemented")

func initLocal() {
	localLoc = *UTC
}

func loadLocation(name string) (*Location, error) {
	return nil, errNotImplemented
}

func readFile(name string) ([]byte, error) {
	return nil, errNotImplemented
}

func open(name string) (uintptr, error) {
	return 0, errNotImplemented
}

func closefd(fd uintptr) {}

func preadn(fd uintptr, buf []byte, off int) error {
	return errNotImplemented
}
