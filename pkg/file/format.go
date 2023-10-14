package file

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

type headerSeeks struct {
	magic      int64
	client     int64
	version    int64
	ntz        int64
	fss        int64
	terminator int64
	data       int64
}

var DefaultHeaderSeeks headerSeeks = headerSeeks{
	H_MAGICBEGIN,
	H_CLIENTBEGIN,
	H_VERSIONBEGIN,
	H_NTZBEGIN,
	H_FSSBEGIN,
	-1,
	-1,
}

type HeaderInfo struct {
	Client  string
	Version string
	Ntz     uint32
	Fss     uint32
}

var HEADER_MAGIC_NUMBER = []byte{1, 9, 6, 5, 0, '\n'}
var HEADER_TERMINATOR = []byte("==== END HEADER ====\n")

func MakeFreeSpaceBuffer(freeSpaceLines []string) ([]byte, uint32, error) {
	buffer := []byte{}

	fss := 0
	for _, s := range freeSpaceLines {
		if strings.Index(s, "\n") != -1 {
			return nil, 0, errors.New("Invalid character '\\n' in free-space line!")
		}
		fss += len(s) + 1
		buffer = append(buffer, []byte(s+"\n")...)
	}

	return buffer, uint32(fss), nil
}

/*
Creates a writable byte buffer of the header and updates info.Fss
*/
func MakeHeader(info HeaderInfo, freeSpace []byte) ([]byte, error) {
	headerBuffer := make([]byte, H_MINSIZE+len(freeSpace))
	copy(headerBuffer[0:H_CLIENTBEGIN], HEADER_MAGIC_NUMBER)
	copy(headerBuffer[H_CLIENTBEGIN:H_DIVIDERBEGIN], []byte(info.Client))
	headerBuffer[H_DIVIDERBEGIN] = '/'
	copy(headerBuffer[H_VERSIONBEGIN:H_NTZBEGIN], []byte(info.Version))

	binary.BigEndian.PutUint32(headerBuffer[H_NTZBEGIN:H_FSSBEGIN], info.Ntz)
	headerBuffer[H_NTZEND-1] = '\n'

	binary.BigEndian.PutUint32(headerBuffer[H_FSSBEGIN:H_FSSEND], info.Fss)
	headerBuffer[H_FSSEND-1] = '\n'
	copy(headerBuffer[H_FREEBEGIN:H_FREEBEGIN+int64(len(freeSpace))], freeSpace)
	copy(headerBuffer[H_FREEBEGIN+int64(len(freeSpace)):H_MINSIZE+len(freeSpace)], HEADER_TERMINATOR)

	return headerBuffer, nil
}

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
