package file

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"

	"github.com/ImFstAsFckBoi/flock/pkg/utils"
)

// A object to write encrypted data to a .locked file
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

// Create a new LockWriter
func NewLockWriter(path string, passwd *[]byte, client string, version string, perms fs.FileMode) (*LockWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, perms)
	if err != nil {
		return nil, err
	}

	salt := [16]byte(make([]byte, 16))
	n, err := rand.Reader.Read(salt[:])
	if err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(*passwd), salt[:], 1, 64*1024, 1, 32)

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	lw := LockWriter{
		f,
		0,
		c,
		nil,
		DefaultHeaderSeeks,
		make([]byte, 16),
		0,
	}

	lw.Info, err = NewHeaderInfo(client, version, salt, blake2b.Sum256(key), []string{"Meow! :3"})

	lw.seeks = lw.Info.GetSeeks()

	header, err := lw.Info.MakeHeader()

	if err != nil {
		return nil, err
	}

	n, err = f.Write(header)

	if n != len(header) || err != nil {
		return nil, errors.New("Failed to make file header")
	}

	return &lw, nil
}

// Write encrypted data to file
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
			_, err := lw.flushBuffer()
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

// Close Writer, flush whats left in buffer
func (lw *LockWriter) Close() (int, error) {
	n, err := lw.flushBuffer()
	err = errors.Join(err, lw.flushHeader())
	err = errors.Join(err, lw.File.Close())

	return n, err
}

// Flush internal buffer to file
func (lw *LockWriter) flushBuffer() (int, error) {
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

// Flush internal in-memory header to file
func (lw *LockWriter) flushHeader() error {
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
