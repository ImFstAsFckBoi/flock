package file

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type HeaderInfo struct {
	Version string
	Ntz     int
}

/* Header format
Magic number:       0x19650
Client/Version:     locker/0.1.2
trailning zeros:    54 (32 bytes for expanding ntz)
Free space(unused): gz_compress=1
                    aes: 256
                    foo=bar
                    xyz123
                    Meow! :3
===== END HEADER =====
data...
*/

var HEADER_MAGIC_NUMBER = "\001\011\006\005\000\n"
var HEADER_VERSION = "locker/%VERSION\n"
var HEADER_NTZ = "%NTZ                                \n"
var HEADER_TERMINATOR = "===== END HEADER =====\n"
var HEADER_FMT = HEADER_MAGIC_NUMBER + HEADER_VERSION + HEADER_NTZ + HEADER_TERMINATOR

func MakeHeader(version string, ntz int) string {
	fmt := HEADER_FMT
	fmt = strings.Replace(fmt, "%VERSION", version, 1)
	fmt = strings.Replace(fmt, "%NTZ", strconv.Itoa(ntz), 1)
	return fmt
}

func HasMagicNumber(data []byte) bool {
	return string(data[0:6]) == HEADER_MAGIC_NUMBER
}

func ReadHeaderInfo(data []byte) (*HeaderInfo, int, error) {

	if !HasMagicNumber(data) {
		return &HeaderInfo{"0", 0}, 0, errors.New("Data does not start with")
	}

	h_end := bytes.Index(data, []byte(HEADER_TERMINATOR))
	if h_end == -1 {
		return &HeaderInfo{"0", 0}, 0, errors.New("Could not find header terminator")
	}
	h_end += len(HEADER_TERMINATOR)

	header := data[0:h_end]

	h_lines := bytes.Split(header, []byte("\n"))

	v := string(h_lines[1][7:])
	ntz, err := strconv.Atoi(strings.Trim(string(h_lines[2]), " "))
	if err != nil {
		return &HeaderInfo{"0", 0}, 0, errors.New("Could not reas corrupt header")
	}

	return &HeaderInfo{v, ntz}, h_end, nil

}
