PKG := github.com/ImFstAsFckBoi/locker

build:
	go build ${PKG}/cmd/locker

run:
	go run ${PKG}/cmd/locker
