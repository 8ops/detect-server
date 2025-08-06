# Makefile for detect-server

# Variables
BINARY=detect-server
SOURCE_FILES=*.go
BUILD_DIR=build

# Default target
all: build

# Build the binary
build:
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY} .

# Install dependencies
deps:
	go mod tidy

# Run tests (placeholder)
test:
	@echo "No tests implemented yet"

# Clean build artifacts
clean:
	rm -rf ${BUILD_DIR}

# Install the binary
install: build
	sudo cp ${BUILD_DIR}/${BINARY} /usr/local/bin/

# Help
help:
	@echo "Available targets:"
	@echo "  all     - Build the project (default)"
	@echo "  build   - Build the binary"
	@echo "  deps    - Install dependencies"
	@echo "  test    - Run tests"
	@echo "  clean   - Clean build artifacts"
	@echo "  install - Install the binary to /usr/local/bin/"
	@echo "  help    - Show this help message"

.PHONY: all build deps test clean install help