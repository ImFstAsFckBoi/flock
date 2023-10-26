package file

// Constant header seek values
const (
	H_MAGICBEGIN   int64 = 0
	H_CLIENTBEGIN  int64 = 6
	H_DIVIDERBEGIN int64 = 36
	H_VERSIONBEGIN int64 = 37
	H_SALTBEGIN    int64 = 68
	H_HASHBEGIN    int64 = 85
	H_NTZBEGIN     int64 = 118
	H_FSSBEGIN     int64 = 123
	H_FREEBEGIN    int64 = 128

	H_MAGICEND   int64 = 6
	H_CLIENTEND  int64 = 36
	H_DIVIDEREND int64 = 37
	H_VERSIONEND int64 = 68
	H_SALTEND    int64 = 85
	H_HASHEND    int64 = 118
	H_NTZEND     int64 = 123
	H_FSSEND     int64 = 128

	H_MAGICLEN   int = 6
	H_CLIENTLEN  int = 30
	H_DIVIDERLEN int = 1
	H_VERSIONLEN int = 31
	H_SALTLEN    int = 17
	H_HASHLEN    int = 33
	H_NTZLEN     int = 5
	H_FSSLEN     int = 5
	H_TERMLEN    int = 21

	H_MINSIZE   int = 149
	H_PREFSSIZE int = 128
)

var HEADER_MAGIC_NUMBER = []byte{1, 9, 6, 5, 0, '\n'}
var HEADER_TERMINATOR = []byte("==== END HEADER ====\n")
