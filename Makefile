.PHONY: build run clean test frontend-build

BINARY_NAME=chronicle

build: frontend-build
	go build -o bin/$(BINARY_NAME) main.go

run:
	go run main.go server

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -rf frontend/dist/

frontend-build:
	cd frontend && npm install && npm run build

help:
	@echo "Usage:"
	@echo "  make build           - Build frontend and backend"
	@echo "  make run             - Run the server"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean builds"
	@echo "  make frontend-build  - Build frontend only"
