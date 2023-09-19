package main

import (
	"crypto/aes"
	"fmt"
	"os"

	"github.com/ImFstAsFckBoi/locker/pkg/encrypt"
	"golang.org/x/term"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Read  file
	data, err := os.ReadFile("./test.txt")
	check(err)

	fmt.Print("Enter Password: ")
	rawPasswd, err := term.ReadPassword(0)
	check(err)
	fmt.Println("")
	key := encrypt.BytePadRepeat(rawPasswd, 32)

	cipher, err := aes.NewCipher(key)
	check(err)

	padSize := encrypt.Ceil128(len(data))

	encBuff := make([]byte, padSize)

	cipher.Encrypt(encBuff, encrypt.BytePadZero(data, padSize))

	fmt.Println(encBuff)

	dencBuff := make([]byte, padSize)
	cipher.Decrypt(dencBuff, encBuff)

	fmt.Println(string(dencBuff))
}
