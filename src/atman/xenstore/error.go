package xenstore

import (
	"fmt"
)

type Error interface {
	error

	Retry() bool
}

type responseError string

func (e responseError) Error() string {
	return string(e)
}

func (e responseError) Retry() bool {
	return e == "EAGAIN"
}

type requestError struct {
	requestType uint32
	requestPath string

	error
}

func (e requestError) Error() string {
	return fmt.Sprintf(
		"op=%d path=%q err=%q",
		e.requestType,
		e.requestPath,
		e.error,
	)
}

func (e requestError) Retry() bool {
	return isRetry(e.error)
}

func isRetry(e error) bool {
	err, ok := e.(Error)
	return ok && err.Retry()
}

func annotateError(requestType uint32, requestPath string, err error) error {
	return &requestError{
		requestType: requestType,
		requestPath: requestPath,
		error:       err,
	}
}
