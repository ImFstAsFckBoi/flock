package utils

import "golang.org/x/exp/constraints"

/*
The Memset() function fills the first n bytes of the memory
area pointed to by s with the constant byte c.
*/
func Memset[T constraints.Integer](dst []T, c T) {
	for i := 0; i < len(dst); i++ {
		dst[i] = c
	}
}
