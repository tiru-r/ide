# GoX IDE Makefile

# Build variables
BINARY_NAME=gox-ide
VERSION=0.1.0-alpha
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: all build clean test fmt vet lint install dev

# Default target
all: clean fmt vet test build

# Build the project
build:
	@echo "🔨 Building GoX IDE..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "✅ Build complete: $(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "🔨 Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	@echo "✅ Cross-platform build complete"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	go clean
	rm -f $(BINARY_NAME)
	rm -rf dist/
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...
	@echo "✅ Tests complete"

# Format code
fmt:
	@echo "📝 Formatting code..."
	go fmt ./...
	@echo "✅ Format complete"

# Run vet
vet:
	@echo "🔍 Running vet..."
	go vet ./...
	@echo "✅ Vet complete"

# Run linter (requires golangci-lint)
lint:
	@echo "🔍 Running linter..."
	golangci-lint run
	@echo "✅ Lint complete"

# Install binary to GOPATH/bin
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	mv $(BINARY_NAME) $(GOPATH)/bin/
	@echo "✅ Install complete"

# Development mode with hot reload
dev:
	@echo "🔄 Starting development mode..."
	go run . .

# Show version info
version:
	@echo "GoX IDE v$(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@echo "Git commit: $(GIT_COMMIT)"

# Setup development environment
setup:
	@echo "🚀 Setting up development environment..."
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Setup complete"

# Release preparation
release: clean fmt vet test build-all
	@echo "🚀 Release $(VERSION) ready!"