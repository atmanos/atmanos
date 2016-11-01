package rand

import (
	"io"
	"math/rand"
)

func init() {
	Reader = newInsecureReader()
}

type insecureReader struct{}

func newInsecureReader() io.Reader {
	return &insecureReader{}
}

func (*insecureReader) Read(buf []byte) (n int, err error) {
	return rand.Read(buf)
}
