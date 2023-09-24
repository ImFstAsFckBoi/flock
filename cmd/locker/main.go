package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"

	"github.com/ImFstAsFckBoi/locker/pkg/encrypt"
	"github.com/ImFstAsFckBoi/locker/pkg/file"
	"github.com/ImFstAsFckBoi/locker/pkg/utils"
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
	data, err := os.ReadFile(path)
	utils.ErrCheck(err)

	padSize := encrypt.Ceil32(len(data))
	utils.ErrCheck(err)

	encBuff := make([]byte, padSize)

	cipher.Encrypt(encBuff, encrypt.BytePadZero(data, padSize))

	header := file.MakeHeader("0.1", padSize-len(data))

	f, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0600)
	utils.ErrCheck(err)
	_, err = f.Write([]byte(header))
	utils.ErrCheck(err)

	_, err = f.Write(encBuff)
	utils.ErrCheck(err)
	f.Close()

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
