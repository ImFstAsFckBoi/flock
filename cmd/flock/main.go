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
	cipher, err := GetPasswordCipher("Enter password: ")
	utils.ErrCheck(err)

	// path, err := GetFileArgs()
	path := "./test.txt"
	utils.ErrCheck(err)

	write := true
	read := true

	if write {
		f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
		utils.ErrCheck(err)
		info, err := f.Stat()
		utils.ErrCheck(err)
		defer f.Close()
		writer, err := file.NewLockWriter(
			path+".locked",
			&cipher,
			"flock",
			"0.1",
			info.Mode().Perm(),
		)

		utils.ErrCheck(err)

		chunkSz := 20

		inBuff := make([]byte, chunkSz)
		// encBuff := make([]byte, chunkSz)
		n := chunkSz
		for n == chunkSz {
			n, err = f.Read(inBuff)
			utils.ErrCheck(err)
			_, err = writer.Write(inBuff[0:n])
			utils.ErrCheck(err)
		}

		writer.Close()
	}

	if read {
		reader, err := file.NewLockReader(path+".locked", &cipher)
		utils.ErrCheck(err)

		println("Free-space:")
		for _, s := range reader.Info.FreeSpace {
			println(s)
		}
		println("\nContent:")

		chunkSz := 8

		inBuff := make([]byte, chunkSz)
		n := chunkSz
		for n == chunkSz {
			n, err = reader.Read(inBuff)
			utils.ErrCheck(err)
			print(string(inBuff))
		}

		reader.Close()
	}
}
