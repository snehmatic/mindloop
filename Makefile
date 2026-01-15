# Binary names
BINARY_NAME=mindloop
SERVER_BINARY_NAME=mindloop-server

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build build-cli build-server run-server test fmt clean help

## all: Build both CLI and Server binaries
all: build

## build: Build both CLI and Server binaries
build: build-cli build-server

## build-cli: Build the CLI binary
build-cli:
	@echo "  >  Building CLI binary..."
	go build -o $(BINARY_NAME) main.go

## build-server: Build the Server binary
build-server:
	@echo "  >  Building Server binary..."
	go build -o $(SERVER_BINARY_NAME) cmd/server/server.go

## run-server: Run the server directly
run-server:
	@echo "  >  Running server..."
	go run cmd/server/server.go

## test: Run all unit tests
test:
	@echo "  >  Running tests..."
	go test ./...

## fmt: Format all go files
fmt:
	@echo "  >  Formatting code..."
	go fmt ./...

## clean: Clean build files
clean:
	@echo "  >  Cleaning build cache..."
	go clean
	rm -f $(BINARY_NAME) $(SERVER_BINARY_NAME)

## help: Show help
help: Makefile
	@echo
	@echo " Choose a command run in "$(AppName)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
