package exec

import "errors"

var ErrNotSupported = errors.New("not supported")

func LookPath(file string) (string, error) {
	return "", ErrNotSupported
}
