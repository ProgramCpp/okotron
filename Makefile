.PHONY: build

all: build test build-deb

build: 
	go build -o ./build/

test:
	go test ./... 

run:
	./build/okotron
	
build-deb:
	env GOOS=debian GOARCH=x86/64 go build -o ./build/oktron-deb
