.PHONY: build run clean test help

BINARY_NAME=chronicle$(shell go env GOEXE)

# Git info
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null)
GIT_DATE=$(shell git log -1 --format=%cs 2>/dev/null)
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null)

LDFLAGS=-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.gitDate=$(GIT_DATE) -X main.buildTime=$(BUILD_TIME)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) main.go

run:
	go run main.go server

test:
	go test -v ./...

clean:
	rm -rf bin/

help:
	@echo "Usage:"
	@echo "  make build           - Build the CLI and API server"
	@echo "  make run             - Run the server"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean builds"
