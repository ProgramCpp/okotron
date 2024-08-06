.PHONY: build

all: build test

build: 
	go build -o ./build/

test:
	go test ./... 

run:
	./build/okotron
	