# Binary names
BINARY_NAME=mindloop
SERVER_BINARY_NAME=mindloop-server

# Defaults for server configuration
PORT ?= 8765
MODE ?= local

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build build-cli build-server run-server start-server kill-server test fmt clean help

## all: Build both CLI and Server binaries
all: build

## build: Build the mindloop binary
build:
	@echo "  >  Building binary..."
	go build -o $(BINARY_NAME) main.go

## run-server: Run the server directly
run-server:
	@echo "  >  Running server on port $(PORT)..."
	go run main.go server -p $(PORT)

## start-server: Build and run the server in background
start-server: build
	@if [ -f .mindloop-server.pid ]; then \
		if ps -p $$(cat .mindloop-server.pid) > /dev/null 2>&1; then \
			echo "  >  Server is already running (PID: $$(cat .mindloop-server.pid))"; \
			exit 1; \
		else \
			echo "  >  Removing stale PID file..."; \
			rm -f .mindloop-server.pid; \
		fi; \
	fi
	@echo "  >  Starting server in background on port $(PORT)..."
	@./$(BINARY_NAME) server -p $(PORT) > server.log 2>&1 & echo $$! > .mindloop-server.pid
	@sleep 0.5
	@if ps -p $$(cat .mindloop-server.pid) > /dev/null 2>&1; then \
		echo "  >  Server started (PID: $$(cat .mindloop-server.pid))"; \
		echo "  >  Logs: server.log"; \
	else \
		echo "  >  Server failed to start. Check server.log for errors."; \
		rm -f .mindloop-server.pid; \
		exit 1; \
	fi

## kill-server: Stop the background server
kill-server:
	@if [ ! -f .mindloop-server.pid ]; then \
		echo "  >  No PID file found. Server may not be running."; \
		exit 1; \
	fi
	@PID=$$(cat .mindloop-server.pid); \
	if ps -p $$PID > /dev/null 2>&1; then \
		echo "  >  Stopping server (PID: $$PID)..."; \
		kill -TERM $$PID; \
		sleep 1; \
		if ps -p $$PID > /dev/null 2>&1; then \
			echo "  >  Server still running, forcing shutdown..."; \
			kill -9 $$PID; \
		fi; \
		echo "  >  Server stopped."; \
	else \
		echo "  >  Server process not found (stale PID: $$PID)."; \
	fi
	@rm -f .mindloop-server.pid

## test: Run all unit tests
test:
	@echo "  >  Running tests..."
	go test ./...

## lint: Run linters
lint:
	@echo "  >  Running linters..."
	golangci-lint run

## fmt: Format all go files
fmt:
	@echo "  >  Formatting code..."
	go fmt ./...

## clean: Clean build files
clean:
	@echo "  >  Cleaning build cache..."
	go clean
	rm -f $(BINARY_NAME) $(SERVER_BINARY_NAME) server.log .mindloop-server.pid mindloop.log

## help: Show help
help: Makefile
	@echo
	@echo " Choose a command run in "$(AppName)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
