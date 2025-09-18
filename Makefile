# Go Build and Test Makefile

.PHONY: build test clean run install lint fmt vet

# Binary name
BINARY_NAME=weblate-mcp

# Go related variables
GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_FILES := $(shell find . -name '*.go' -type f)

# Build the binary
build:
	go build -ldflags '-w -s' -o $(BINARY_NAME) .

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags '-w -s' -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o $(BINARY_NAME)-windows-amd64.exe .

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)*
	rm -f coverage.out coverage.html

# Run the application
run:
	go run .

# Install dependencies
install:
	go mod download
	go mod tidy

# Lint the code
lint:
	golangci-lint run

# Format the code
fmt:
	go fmt ./...

# Vet the code
vet:
	go vet ./...

# Check all (fmt, vet, test)
check: fmt vet test

# Development setup
dev-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Build for production
prod-build: clean fmt vet test build

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Run the application"
	@echo "  install       - Install dependencies"
	@echo "  lint          - Lint the code"
	@echo "  fmt           - Format the code"
	@echo "  vet           - Vet the code"
	@echo "  check         - Run fmt, vet, and test"
	@echo "  dev-setup     - Install development tools"
	@echo "  prod-build    - Production build (clean, fmt, vet, test, build)"
	@echo "  help          - Show this help"