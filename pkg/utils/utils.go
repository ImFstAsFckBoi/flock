package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ErrCheck(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	BIAS_YES int = iota
	BIAS_NO
	BIAS_NONE
)

func YNPrompt(msg string, bias int) (bool, error) {
	y := "y"
	n := "n"

	switch bias {
	case BIAS_YES:
		y = strings.ToUpper(y)
	case BIAS_NO:
		n = strings.ToUpper(n)
	case BIAS_NONE:
	default:
		return false, errors.New("Invalid bias argument. Must be BIAS_Y, BIAS_N or BIAS_NONE")
	}

	res := false

	stdin := bufio.NewReader(os.Stdin)
	for true {
		stdin.Discard(stdin.Buffered())
		fmt.Printf("%s [%s/%s]", msg, y, n)
		r, len, err := stdin.ReadRune()
		if err != nil || len != 1 {
			return false, err
		}
		switch r {
		case 'y', 'Y':
			res = true
		case 'n', 'N':
			res = false
		case '\n', '\r':
			switch bias {
			case BIAS_YES:
				res = true
			case BIAS_NO:
				res = false
			default:
				continue
			}
		default:
			continue
		}

		break
	}

	return res, nil
}
