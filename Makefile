.PHONY: build

all: build test build-deb

build: 
	go build -o ./build/

test:
	go test ./... 

run:
	./build/okotron
	
# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
build-deb:
	env GOOS=linux GOARCH=amd64 go build -o ./build/okotron-deb
