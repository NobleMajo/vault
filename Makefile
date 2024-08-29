BINARY_NAME := vault
VERSION := $(shell cat version.txt)
COMMIT := $(shell git rev-parse --short HEAD)

build:
	go build -o $(BINARY_NAME) -ldflags='-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)' main.go

run:
	go run main.go

clean:
	go clean
	rm -f $(BINARY_NAME)

deps:
	go mod download

update:
	rm -f go.sum
	go mod tidy
	go get -t -u ./...

test:
	go test ./...

coverage:
	go test ./... -cover