package file

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"os"

	"github.com/ImFstAsFckBoi/flock/pkg/utils"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

type LockReader struct {
	File        *os.File
	cipher      cipher.Block
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
	16,
	16,
}

func NewLockReader(path string, passwd *[]byte) (*LockReader, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	rw := LockReader{
		f,
		nil,
		HeaderInfo{},
		make([]byte, 16),
		16,
		16,
	}

	err = rw.Info.ReadHeader(rw.File)
	if err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(*passwd), rw.Info.Salt[:], 1, 64*1024, 1, 32)

	rw.cipher, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if rw.Info.Hash != blake2b.Sum256(key) {
		return nil, errors.New("Password key does not match hash in header!")
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
			_, err := rw.FillBuffer()
			if err == io.EOF && bIdx == 0 {
				return bIdx, err
			} else if err == io.EOF {
				utils.Memset[byte](b[bIdx:bLen], 0)
				return bIdx, nil
			}
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

func (rw *LockReader) FillBuffer() (int, error) {
	// TODO: Deal with n | 16 edge case
	n, err := rw.File.Read(rw.buffer)
	if err != nil {
		return n, err
	}

	rw.cipher.Decrypt(rw.buffer, rw.buffer)

	// Test if EOF is next
	seek, err1 := rw.File.Seek(0, io.SeekCurrent)
	info, err2 := rw.File.Stat()

	if errors.Join(err1, err2) != nil {
		return n, err
	}

	if seek == info.Size() {
		rw.bufferEnd = 16 - int(rw.Info.Ntz)
	}

	rw.bufferStart = 0

	return n, nil
}

func (rw *LockReader) Close() error {
	return rw.File.Close()
}
