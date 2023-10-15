package file

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/ImFstAsFckBoi/flock/pkg/utils"
)

type LockWriter struct {
	File        *os.File
	Count       int
	cipher      cipher.Block
	Info        *HeaderInfo
	seeks       HeaderSeeks
	buffer      []byte
	bufferCount int
}

var deadWriter LockWriter = LockWriter{
	nil,
	0,
	nil,
	nil,
	HeaderSeeks{},
	nil,
	0,
}

/*
Wrapper around file to automatically write the lock file header and encrypt
all written in chunks.
*/
func NewLockWriter(path string, cipher cipher.Block, client string, version string, perms fs.FileMode) (*LockWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, perms)
	if err != nil {
		return nil, err
	}

	lw := LockWriter{
		f,
		0,
		cipher,
		nil,
		DefaultHeaderSeeks,
		make([]byte, 16),
		0,
	}

	lw.Info, err = NewHeaderInfo(client, version, []string{"Meow! :3"})

	lw.seeks = lw.Info.GetSeeks()

	header, err := lw.Info.MakeHeader()

	if err != nil {
		return nil, err
	}

	n, err := f.Write(header)

	if n != len(header) || err != nil {
		return nil, errors.New("Failed to make file header")
	}

	return &lw, nil
}

func (lw *LockWriter) Write(p []byte) (int, error) {
	pLen := len(p)

	// TODO: REMOVE DEBUG
	// fmt.Printf("Write content: %s\n", string(p))

	pIdx := 0
	flushes := 0
	for pIdx < pLen {
		needed := 16 - lw.bufferCount
		have := pLen - pIdx

		if needed < have {
			copy(lw.buffer[lw.bufferCount:16], p[pIdx:pIdx+needed])

			pIdx += needed
			lw.bufferCount += needed
		} else {
			copy(lw.buffer[lw.bufferCount:lw.bufferCount+have], p[pIdx:pLen])

			pIdx += have
			lw.bufferCount += have
		}

		if lw.bufferCount >= 16 {
			_, err := lw.FlushBuffer()
			if err != nil {
				return flushes * 16, err
			}

			flushes += 1
		}

	}
	// TODO: REMOVE DEBUG
	// fmt.Printf("Buffer content on exit: %s\n", string(lw.buffer))

	return flushes * 16, nil
}

func (lw *LockWriter) Close() (int, error) {
	n, err := lw.FlushBuffer()
	err = errors.Join(err, lw.FlushHeader())
	err = errors.Join(err, lw.File.Close())

	return n, err
}

func (lw *LockWriter) FlushBuffer() (int, error) {
	if lw.bufferCount == 0 {
		return 0, nil
	}

	// TODO: REMOVE DEBUG
	// fmt.Printf("Buffer content: %s\n", string(lw.buffer))

	lw.cipher.Encrypt(lw.buffer, lw.buffer)

	n, err := lw.File.Write(lw.buffer)

	if err != nil {
		return n, errors.New(
			fmt.Sprintf("Failed to write buffer to '%s'", lw.File.Name()),
		)
	}

	lw.Info.Ntz = uint32(16 - lw.bufferCount)
	lw.bufferCount = 0

	utils.Memset[byte](lw.buffer, 0)

	return n, nil
}

func (lw *LockWriter) FlushHeader() error {
	if len(lw.Info.Client) > 29 {
		return errors.New("Headers 'client' field exceeds maximum length of 29 characters")
	} else if len(lw.Info.Version) > 29 {
		return errors.New("Headers 'version' field exceeds maximum length of 29 characters")
	}

	lw.File.WriteAt([]byte(lw.Info.Client), lw.seeks.client)
	lw.File.WriteAt([]byte(lw.Info.Version), lw.seeks.version)

	ntzBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(ntzBuffer, lw.Info.Ntz)
	lw.File.WriteAt(ntzBuffer, lw.seeks.ntz)
	return nil
}
