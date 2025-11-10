# AutoNode Makefile
# Provides commands for building cross-platform binaries

# Binary name
BINARY_NAME=autonode

# Version (can be overridden: make VERSION=1.0.0 build)
VERSION?=0.5.0

# Build directory
BUILD_DIR=bin

# Source directory
CMD_DIR=./cmd/autonode

# Go build flags
LDFLAGS=-ldflags "-X main.version=${VERSION}"

.PHONY: all build clean test install linux windows darwin help

# Default target
all: build

# Build for current platform
build:
	@echo "Building ${BINARY_NAME} for current platform..."
	@mkdir -p ${BUILD_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${CMD_DIR}
	@echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"

# Build for all platforms
build-all: linux windows darwin
	@echo "All platform builds complete!"

# Build for Linux (amd64 and arm64)
linux:
	@echo "Building for Linux amd64..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ${CMD_DIR}
	@echo "Building for Linux arm64..."
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 ${CMD_DIR}

# Build for macOS (amd64 and arm64/Apple Silicon)
darwin:
	@echo "Building for macOS amd64 (Intel)..."
	@mkdir -p ${BUILD_DIR}
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 ${CMD_DIR}
	@echo "Building for macOS arm64 (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 ${CMD_DIR}

# Build for Windows (amd64)
windows:
	@echo "Building for Windows amd64..."
	@mkdir -p ${BUILD_DIR}
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ${CMD_DIR}

# Install binary to system (requires sudo on Unix systems)
install: build
	@echo "Installing ${BINARY_NAME} to /usr/local/bin..."
	@sudo cp ${BUILD_DIR}/${BINARY_NAME} /usr/local/bin/
	@echo "Installation complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ${BUILD_DIR}
	@go clean
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Help
help:
	@echo "AutoNode Makefile Commands:"
	@echo ""
	@echo "  make build        - Build for current platform"
	@echo "  make build-all    - Build for all platforms (Linux, macOS, Windows)"
	@echo "  make linux        - Build for Linux (amd64, arm64)"
	@echo "  make darwin       - Build for macOS (amd64, arm64)"
	@echo "  make windows      - Build for Windows (amd64)"
	@echo "  make install      - Install binary to /usr/local/bin"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make tidy         - Tidy dependencies"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=${VERSION}  - Override with: make VERSION=x.y.z build"
