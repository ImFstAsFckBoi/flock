package utils

func ErrCheck(err error) {
	if err != nil {
		panic(err)
	}
}
