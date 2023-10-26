package file

import (
	"bytes"
)

// All seeks in a file header
type HeaderSeeks struct {
	magic      int64
	client     int64
	version    int64
	ntz        int64
	fss        int64
	freeSpace  int64
	terminator int64
	data       int64
}

var DefaultHeaderSeeks HeaderSeeks = HeaderSeeks{
	H_MAGICBEGIN,
	H_CLIENTBEGIN,
	H_VERSIONBEGIN,
	H_NTZBEGIN,
	H_FSSBEGIN,
	H_FSSEND,
	-1,
	-1,
}

// Check is data array contains magic number
func IsMagicNumber(data []byte) bool {
	return bytes.Compare(data, HEADER_MAGIC_NUMBER) == 0
}
