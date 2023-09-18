package main

import (
	"crypto/aes"
	"fmt"
	"math"
	"os"

	"golang.org/x/term"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Round up to nearest multiple of 128
func Ceil128(x int) int {
	return 128 * int(math.Ceil(float64(x)/float64(128)))
}

// Will pad, (or slice), the b to fit length n. Padded space will be filled by repeating b.
func BytePadRepeat(b []byte, n int) []byte {
	out := b
	if len(b) >= n {
		out = b[0:n]
	} else {
		p_len := len(b)
		rem_len := n - p_len

		for rem_len != 0 {
			add_len := min(rem_len, p_len)
			out = append(out, b[0:add_len]...)
			rem_len -= add_len
		}
	}

	return out
}

func BytePadZero(b []byte, n int) []byte {
	out := b
	diff := n - len(b)

	if diff < 0 {
		out = b[0:n]
	} else {
		out = append(out, make([]byte, diff)...)
	}

	return out
}

func main() {
	// Read  file
	data, err := os.ReadFile("./test.txt")
	check(err)
	fmt.Println(string(data))

	fmt.Print("Enter Password: ")
	rawPasswd, err := term.ReadPassword(0)

	key := BytePadRepeat(rawPasswd, 32)

	cipher, err := aes.NewCipher(key)
	check(err)

	padSize := Ceil128(len(data))

	encBuff := make([]byte, padSize)

	cipher.Encrypt(encBuff, BytePadZero(data, padSize))

	fmt.Println(string(encBuff))

	dencBuff := make([]byte, padSize)
	cipher.Decrypt(dencBuff, encBuff)

	fmt.Println(string(dencBuff))
}
