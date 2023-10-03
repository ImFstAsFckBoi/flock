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
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	writer, err := file.NewLockWriter(path+".lock", "0.1")
	utils.ErrCheck(err)

	chunkSz := 32

	inBuff := make([]byte, chunkSz)
	encBuff := make([]byte, chunkSz)
	n := chunkSz
	for n == chunkSz {
		n, err = f.Read(inBuff)
		utils.ErrCheck(err)
		cipher.Encrypt(encBuff, encrypt.BytePadZero(inBuff, chunkSz))
		_, err = writer.Write(encBuff)
		utils.ErrCheck(err)
	}

	err = writer.FlushHeader()
	// Read .lock

	padSize := encrypt.Ceil32(writer.WriteCount)

	fileData, err := os.ReadFile(path + ".lock")

	info, dataStart, err := file.ReadHeaderInfo(fileData)

	encData := fileData[dataStart:]

	dencBuff := make([]byte, padSize)
	cipher.Decrypt(dencBuff, encData)

	dencBuff = dencBuff[0 : len(dencBuff)-info.Ntz]

	f, err = os.OpenFile(path+".unlocked", os.O_CREATE|os.O_RDWR, 0600)
	utils.ErrCheck(err)
	_, err = f.Write(dencBuff)
	utils.ErrCheck(err)
	f.Close()
}
