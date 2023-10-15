package file

import (
	"encoding/binary"
	"errors"
	"strings"
)

type HeaderInfo struct {
	Client    string
	Version   string
	Salt      [16]byte
	Hash      [32]byte
	Ntz       uint32
	Fss       uint32
	FreeSpace []string
}

func NewHeaderInfo(client string, version string, salt [16]byte, hash [32]byte, freeSpace []string) (*HeaderInfo, error) {
	info := HeaderInfo{client, version, salt, hash, 0, 0, nil}
	err := info.UpdateFS(freeSpace)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (info *HeaderInfo) UpdateFSS() uint32 {
	info.Fss = 0
	for _, s := range info.FreeSpace {
		info.Fss += uint32(len(s) + 1)
	}
	return info.Fss
}

func (info *HeaderInfo) UpdateFS(freeSpace []string) error {
	info.FreeSpace = freeSpace
	for _, s := range info.FreeSpace {
		if strings.Index(s, "\n") != -1 {
			return errors.New("Invalid character '\\n' in free-space line!")
		}
	}
	info.UpdateFSS()
	return nil
}

func (info *HeaderInfo) FSAsBytes() []byte {
	buffer := make([]byte, info.Fss)
	if info.Fss == 0 || len(info.FreeSpace) == 0 {
		return nil
	}

	idx := 0
	for _, s := range info.FreeSpace {
		copy(buffer[idx:len(s)], []byte(s+"\n"))
		idx += len(s)
	}
	return buffer
}

func (info *HeaderInfo) GetSeeks() HeaderSeeks {
	seeks := DefaultHeaderSeeks
	seeks.terminator = H_FREEBEGIN + int64(info.Fss)
	seeks.data = int64(H_MINSIZE + int(info.Fss))

	return seeks
}

func (info *HeaderInfo) MakeHeader() ([]byte, error) {
	headerBuffer := make([]byte, H_MINSIZE+int(info.Fss))
	bytesFS := info.FSAsBytes()

	copy(headerBuffer[0:H_CLIENTBEGIN], HEADER_MAGIC_NUMBER)
	copy(headerBuffer[H_CLIENTBEGIN:H_CLIENTEND], []byte(info.Client))
	headerBuffer[H_DIVIDERBEGIN] = '/'
	copy(headerBuffer[H_VERSIONBEGIN:H_VERSIONEND], []byte(info.Version))
	headerBuffer[H_VERSIONEND-1] = '\n'

	copy(headerBuffer[H_SALTBEGIN:H_SALTEND], info.Salt[:])
	headerBuffer[H_SALTEND-1] = '\n'

	copy(headerBuffer[H_HASHBEGIN:H_HASHEND], info.Hash[:])
	headerBuffer[H_HASHEND-1] = '\n'

	binary.BigEndian.PutUint32(headerBuffer[H_NTZBEGIN:H_NTZEND], info.Ntz)
	headerBuffer[H_NTZEND-1] = '\n'

	binary.BigEndian.PutUint32(headerBuffer[H_FSSBEGIN:H_FSSEND], info.Fss)
	headerBuffer[H_FSSEND-1] = '\n'
	if bytesFS != nil {
		copy(headerBuffer[H_FREEBEGIN:H_FREEBEGIN+int64(info.Fss)], bytesFS)
		copy(headerBuffer[H_FREEBEGIN+int64(info.Fss):H_MINSIZE+int(info.Fss)], HEADER_TERMINATOR)
	}

	return headerBuffer, nil
}
