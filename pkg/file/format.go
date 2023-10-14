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
	terminator int64
	data       int64
}

var DefaultHeaderSeeks HeaderSeeks = HeaderSeeks{
	H_MAGICBEGIN,
	H_CLIENTBEGIN,
	H_VERSIONBEGIN,
	H_NTZBEGIN,
	H_FSSBEGIN,
	-1,
	-1,
}

var HEADER_MAGIC_NUMBER = []byte{1, 9, 6, 5, 0, '\n'}
var HEADER_TERMINATOR = []byte("==== END HEADER ====\n")

/*
Creates a writable byte buffer of the header and updates info.Fss
*/

func HasMagicNumber(data []byte) bool {
	return bytes.Compare(data, HEADER_MAGIC_NUMBER) == 0
}

//func ReadHeaderInfo(data []byte) (*HeaderInfo, int, error) {

//	if !HasMagicNumber(data) {
//		return &HeaderInfo{}, 0, errors.New("Malformed header!")
//	}

//	h_end := bytes.Index(data, []byte(HEADER_TERMINATOR))
//	if h_end == -1 {
//		return &HeaderInfo{}, 0, errors.New("Could not find header terminator")
//	}
//	h_end += len(HEADER_TERMINATOR)

//	header := data[0:h_end]

//	h_lines := bytes.Split(header, []byte("\n"))

//	v := string(h_lines[1][7:])
//	ntz, err := strconv.Atoi(strings.Trim(string(h_lines[2]), " "))
//	if err != nil {
//		return &HeaderInfo{}, 0, errors.New("Could not reas corrupt header")
//	}

//	// TODO: read client
//	return &HeaderInfo{"adas", v, uint32(ntz), uint32(fss)}, h_end, nil

//}
