package file

const (
	H_MAGICBEGIN   int64 = 0
	H_CLIENTBEGIN  int64 = 6
	H_DIVIDERBEGIN int64 = 36
	H_VERSIONBEGIN int64 = 37
	H_NTZBEGIN     int64 = 68
	H_FSSBEGIN     int64 = 73
	H_FREEBEGIN    int64 = 78

	H_MAGICEND   int64 = 6
	H_CLIENTEND  int64 = 36
	H_DIVIDEREND int64 = 37
	H_VERSIONEND int64 = 68
	H_NTZEND     int64 = 73
	H_FSSEND     int64 = 78

	H_MAGICLEN   int = 6
	H_CLIENTLEN  int = 30
	H_DIVIDERLEN int = 1
	H_VERSIONLEN int = 30
	H_NTZLEN     int = 5
	H_FSSLEN     int = 5
	H_TERMLEN    int = 21

	H_MINSIZE   int = 99
	H_PREFSSIZE int = 78
)
