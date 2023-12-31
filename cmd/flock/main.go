package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/ImFstAsFckBoi/flock/pkg/file"
	"github.com/ImFstAsFckBoi/flock/pkg/utils"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/term"
)

func ReadPassword(prompt string) ([]byte, error) {

	fmt.Print(prompt)
	passwd, err := term.ReadPassword(0)
	fmt.Println("")
	passwdHash := blake2b.Sum256(passwd)

	// Zero out memory for safety ?
	// Not sure if necessary but feels "safter"
	_, err = rand.Reader.Read(passwd)
	if err != nil {
		return nil, err
	}

	return passwdHash[:], err
}

func main() {
	path := "./test.txt"

	write := true
	read := true

	if write {
		passwd, err := ReadPassword("Enter password: ")
		utils.ErrCheck(err)
		f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
		utils.ErrCheck(err)
		info, err := f.Stat()
		utils.ErrCheck(err)
		defer f.Close()
		writer, err := file.NewLockWriter(
			path+".locked",
			&passwd,
			"flock",
			"0.1",
			info.Mode().Perm(),
		)

		utils.ErrCheck(err)

		chunkSz := 100

		inBuff := make([]byte, chunkSz)
		// encBuff := make([]byte, chunkSz)
		n := chunkSz
		wc := 0
		rc := 0
		for n == chunkSz {
			n, err = f.Read(inBuff)
			utils.ErrCheck(err)
			n1, err := writer.Write(inBuff[0:n])
			utils.ErrCheck(err)
			rc += n
			wc += n1
		}

		n, _ = writer.Close()
		wc += n

		fmt.Printf("Read: %d bytes.\n", rc)
		fmt.Printf("Wrote: %d bytes.\n", wc)
	}

	if read {
		passwd, err := ReadPassword("Enter password: ")
		utils.ErrCheck(err)
		reader, err := file.NewLockReader(path+".locked", &passwd)
		utils.ErrCheck(err)

		println("Free-space:")
		for _, s := range reader.Info.FreeSpace {
			println(s)
		}
		println("\nContent:")

		chunkSz := 100

		inBuff := make([]byte, chunkSz)
		n := chunkSz
		rc := 0
		for n == chunkSz {
			n, err = reader.Read(inBuff)
			utils.ErrCheck(err)
			print(string(inBuff))
			rc += n
		}

		fmt.Printf("Read: %d bytes.\n", rc)
		reader.Close()
	}
}
