// The runtime package uses //go:linkname to push a few functions into this
// package but we still need a .s file so the Go tool does not pass -complete
// to the go tool compile so the latter does not complain about Go functions
// with no bodies.
