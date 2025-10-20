.PHONY: help build test test-verbose test-coverage clean run install lint

# Default target
help:
	@echo "Available targets:"
	@echo "  build         - Build the application binary"
	@echo "  test          - Run all unit tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Remove build artifacts"
	@echo "  run           - Build and run the application"
	@echo "  install       - Install dependencies"
	@echo "  lint          - Run go vet and gofmt"

# Build the application
build:
	@echo "Building api-tester..."
	@go build -o api-tester ./cmd/api-tester
	@echo "Build complete: ./api-tester"

# Run all tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f api-tester
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Build and run the application
run: build
	@echo "Running api-tester..."
	@./api-tester

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Lint the code
lint:
	@echo "Running linters..."
	@go vet ./...
	@gofmt -l -w .
	@echo "Linting complete"
