# Variables
BINARY_NAME=anchor
BUILD_DIR=bin

# Targets
.PHONY: all build test clean help

all: build

## build: Build the binary
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-race: Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./...

## Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

## Generate coverage report
test-coverage:
	@echo "Generating coverage..."
	@go test -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "✓ Coverage: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	@# Optional: go clean if you want to clear the build cache
	go clean

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
