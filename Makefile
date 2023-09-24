PKG := github.com/ImFstAsFckBoi/locker
BIN := ./target
BINARY := locker
ENTRY := ${PKG}/cmd/locker

build:
	go build -o ${BIN}/${BINARY} ${PKG}/cmd/locker

build-linux32:
	GOARCH=386   GOOS=linux   go build -o ${BIN}/${BINARY}-linux32 ${ENTRY}

build-linux64:
	GOARCH=386   GOOS=windows go build -o ${BIN}/${BINARY}-win32 ${ENTRY}

build-win32:
	GOARCH=amd64 GOOS=linux   go build -o ${BIN}/${BINARY}-linux64 ${ENTRY}

build-win64:
	GOARCH=amd64 GOOS=windows go build -o ${BIN}/${BINARY}-win64 ${ENTRY}

build-all: build-linux32 build-linux64 build-win32 build-win64

clean:
	rm -rf ${BIN}

run:
	go run ${ENTRY}