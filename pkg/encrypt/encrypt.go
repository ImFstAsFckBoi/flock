package encrypt

import "math"

// Round up to nearest multiple of 32 bytes
func Ceil32(x int) int {
	return 32 * int(math.Ceil(float64(x)/float64(128)))
}

// Pad (or slice) the byte array b to fit length n. Padded space will be filled by repeating b.
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

// Pad (or slice) the byte array b to fit length n. Padded space will be filled by zeros
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
