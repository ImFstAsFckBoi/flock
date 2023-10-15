package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"

	"github.com/ImFstAsFckBoi/flock/pkg/encrypt"
	"github.com/ImFstAsFckBoi/flock/pkg/file"
	"github.com/ImFstAsFckBoi/flock/pkg/utils"
	"golang.org/x/term"
)

func GetPasswordCipher(prompt string) (cipher.Block, error) {

	fmt.Print(prompt)
	rawPasswd, err := term.ReadPassword(0)
	fmt.Println("")
	if err != nil {
		return nil, err
	}

	key := encrypt.BytePadRepeat(rawPasswd, 32)

	return aes.NewCipher(key)
}

func GetFileArgs() (string, error) {
	args := os.Args[1:]
	file := args[0]
	_, err := os.Stat(file)

	if err != nil {
		return "", err
	}

	return file, nil

}

func main() {

	// path, err := GetFileArgs()
	path := "./test.txt"

	write := true
	read := true

	if write {
		cipher, err := GetPasswordCipher("Enter password: ")
		utils.ErrCheck(err)
		f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
		utils.ErrCheck(err)
		info, err := f.Stat()
		utils.ErrCheck(err)
		defer f.Close()
		writer, err := file.NewLockWriter(
			path+".locked",
			cipher,
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
		println("Encrypting...")
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
		cipher, err := GetPasswordCipher("Enter password: ")
		utils.ErrCheck(err)
		reader, err := file.NewLockReader(path+".locked", cipher)
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
