.PHONY: build run test test-race bench lint clean help

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"
BINARY := xxxclaw
BUILD_DIR := bin

## help: Show this help
help:
	@echo "XXXCLAW - Modular Agent System"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build the binary"
	@echo "  make run         Build and run"
	@echo "  make test        Run tests"
	@echo "  make test-race   Run tests with race detector"
	@echo "  make bench       Run benchmarks"
	@echo "  make lint        Run linter"
	@echo "  make clean       Clean build artifacts"
	@echo ""

## build: Compile the binary
build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/xxxclaw/

## run: Build and run the server
run: build
	./$(BUILD_DIR)/$(BINARY)

## test: Run all tests
test:
	go test ./... -v -count=1

## test-race: Run all tests with race detector
test-race:
	go test ./... -race -v -count=1

## bench: Run benchmarks
bench:
	go test ./... -bench=. -benchmem -run=^$

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean -cache -testcache
