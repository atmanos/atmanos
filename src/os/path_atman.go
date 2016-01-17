package os

const (
	PathSeparator     = '/'    // OS-specific path separator
	PathListSeparator = '\000' // OS-specific path list separator
)

// IsPathSeparator reports whether c is a directory separator character.
func IsPathSeparator(c uint8) bool {
	return PathSeparator == c
}
