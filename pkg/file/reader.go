package file

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ImFstAsFckBoi/flock/pkg/utils"
)

type LockReader struct {
	File        *os.File
	cipher      *cipher.Block
	Info        HeaderInfo
	buffer      []byte
	bufferStart int
	bufferEnd   int
}

var deadReader LockReader = LockReader{
	nil,
	nil,
	HeaderInfo{},
	nil,
	32,
	32,
}

func NewLockReader(path string, cipher *cipher.Block) (*LockReader, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	rw := LockReader{
		f,
		cipher,
		HeaderInfo{},
		make([]byte, 32),
		32,
		32,
	}

	err = rw.ReadHeaderInfo()
	if err != nil {
		return nil, err
	}

	return &rw, nil
}

func (rw *LockReader) Read(b []byte) (int, error) {
	dataStart := rw.Info.GetSeeks().data
	curSeek, err := rw.File.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	if curSeek < dataStart {
		rw.File.Seek(dataStart, io.SeekStart)
	}

	bLen := len(b)
	bIdx := 0
	for bIdx < bLen {
		if rw.bufferStart >= rw.bufferEnd {
			// TODO: Replace with custom method
			_, err1 := rw.File.Read(rw.buffer)
			seek, err2 := rw.File.Seek(0, io.SeekCurrent)
			info, err3 := rw.File.Stat()
			if errors.Join(err1, err2, err3) != nil {
				utils.Memset[byte](b[bIdx:bLen], 0)
				return bIdx, err
			}

			if seek == info.Size() {
				rw.bufferEnd = 32 - int(rw.Info.Ntz)
			}

			rw.bufferStart = 0
		}

		needed := bLen - bIdx
		have := rw.bufferEnd - rw.bufferStart

		if needed < have {
			copy(b[bIdx:bLen], rw.buffer[rw.bufferStart:rw.bufferStart+needed])

			bIdx += needed
			rw.bufferStart += needed
		} else {
			copy(b[bIdx:bIdx+have], rw.buffer[rw.bufferStart:rw.bufferEnd])

			bIdx += have
			rw.bufferStart += have
		}
	}

	return bIdx, nil
}

func (rw *LockReader) ReadHeaderInfo() error {
	preBuff := make([]byte, H_PREFSSIZE)
	_, err := rw.File.ReadAt(preBuff, 0)

	if err != nil {
		return err
	}

	if !IsMagicNumber(preBuff[H_MAGICBEGIN:H_MAGICEND]) {
		return errors.New("Malformed header!")
	}
	rw.Info.Client = strings.TrimRight(
		string(preBuff[H_CLIENTBEGIN:H_CLIENTEND]),
		"\000",
	)

	rw.Info.Version = strings.TrimRight(
		string(preBuff[H_VERSIONBEGIN:H_VERSIONEND]),
		"\000",
	)

	rw.Info.Ntz = binary.BigEndian.Uint32(preBuff[H_NTZBEGIN : H_NTZEND-1])
	rw.Info.Fss = binary.BigEndian.Uint32(preBuff[H_FSSBEGIN : H_FSSEND-1])
	if rw.Info.Fss != 0 {

		fsBuff := make([]byte, rw.Info.Fss)
		rw.File.ReadAt(fsBuff, H_FSSEND)
		lines := bytes.Split(fsBuff, []byte{'\n'})
		for _, l := range lines {
			rw.Info.FreeSpace = append(rw.Info.FreeSpace, strings.TrimLeft(string(l), "\000"))
		}
	} else {
		rw.Info.FreeSpace = nil
	}

	return nil
}

func (rw *LockReader) Close() error {
	return rw.File.Close()
}
