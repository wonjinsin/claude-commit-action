# Makefile for cleanarch Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Binary name and paths
BINARY_NAME=server
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/server
BUILD_DIR=./build

# Default target
.PHONY: all
all: clean build

# Build the binary
.PHONY: build
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)

# Run the application
.PHONY: run
run:
	@echo "Running the application..."
	$(GOCMD) run $(MAIN_PATH)/main.go

# Run the application with build
.PHONY: run-build
run-build: build
	@echo "Running built binary..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Test all packages
.PHONY: test
test:
	@echo "Testing..."
	$(GOTEST) -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Format all Go files
.PHONY: fmt
fmt:
	@echo "Formatting..."
	$(GOFMT) -s -w .

# Vet examines Go source code
.PHONY: vet
vet:
	@echo "Vetting..."
	$(GOVET) ./...

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(MAIN_PATH)

# Development mode (with file watching - requires 'air' tool)
.PHONY: dev
dev:
	@if command -v air > /dev/null; then \
		echo "Starting development server with hot reload..."; \
		air; \
	else \
		echo "Air not found. Install it with: go install github.com/cosmtrek/air@latest"; \
		echo "Running normal server..."; \
		$(MAKE) run; \
	fi

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	$(GOCMD) install github.com/cosmtrek/air@latest

# Check if everything is ready to commit
.PHONY: check
check: fmt vet test
	@echo "All checks passed!"

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  all           - Clean and build the project"
	@echo "  build         - Build the binary"
	@echo "  run           - Run the application directly with 'go run'"
	@echo "  run-build     - Build and run the binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format all Go files"
	@echo "  vet           - Vet examines Go source code"
	@echo "  deps          - Download dependencies"
	@echo "  tidy          - Tidy dependencies"
	@echo "  build-linux   - Build for Linux"
	@echo "  dev           - Run in development mode with hot reload (requires air)"
	@echo "  install-tools - Install development tools"
	@echo "  check         - Run fmt, vet, and test"
	@echo "  help          - Show this help message" 