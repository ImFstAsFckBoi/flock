package file

import (
	"bytes"
)

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

var HEADER_MAGIC_NUMBER = []byte{1, 9, 6, 5, 0, '\n'}
var HEADER_TERMINATOR = []byte("==== END HEADER ====\n")

/*
Creates a writable byte buffer of the header and updates info.Fss
*/

func IsMagicNumber(data []byte) bool {
	return bytes.Compare(data, HEADER_MAGIC_NUMBER) == 0
}
