all: build

build: **/*.go
	go build

install: dojo
	cp dojo /usr/local/bin/dojo

clean:
	go clean