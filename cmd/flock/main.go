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

	// Read  file
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
	utils.ErrCheck(err)
	info, err := f.Stat()
	utils.ErrCheck(err)
	defer f.Close()
	writer, err := file.NewLockWriter(
		path+".lock",
		&cipher,
		"flock",
		"0.1",
		info.Mode().Perm(),
	)
	defer writer.Close()

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

	// Read .lock

	//padSize := encrypt.Ceil32(writer.WriteCount)

	//fileData, err := os.ReadFile(path + ".lock")

	//info, dataStart, err := file.ReadHeaderInfo(fileData)

	//encData := fileData[dataStart:]

	//dencBuff := make([]byte, padSize)
	//cipher.Decrypt(dencBuff, encData)

	//dencBuff = dencBuff[0 : len(dencBuff)-info.Ntz]

	//f, err = os.OpenFile(path+".unlocked", os.O_CREATE|os.O_RDWR, 0600)
	//utils.ErrCheck(err)
	//_, err = f.Write(dencBuff)
	//utils.ErrCheck(err)
	//f.Close()
}
