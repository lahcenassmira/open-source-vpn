.PHONY: build clean test test-coverage benchmark install docker-build help

# Build variables
BINARY_SERVER=bin/vpn-server
BINARY_CLIENT=bin/vpn-client
BUILD_FLAGS=-ldflags="-s -w"
GO=go

# Default target
all: build

## build: Build server and client binaries
build:
	@echo "Building VPN server..."
	@mkdir -p bin
	$(GO) build $(BUILD_FLAGS) -o $(BINARY_SERVER) ./cmd/server
	@echo "Building VPN client..."
	$(GO) build $(BUILD_FLAGS) -o $(BINARY_CLIENT) ./cmd/client
	@echo "Build complete!"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf coverage/
	@$(GO) clean
	@echo "Clean complete!"

## test: Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v -race ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p coverage
	$(GO) test -v -race -coverprofile=coverage/coverage.out ./...
	$(GO) tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report generated: coverage/coverage.html"

## benchmark: Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

## test-integration: Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GO) test -v -tags=integration ./tests/integration/...

## install: Install binaries to GOPATH/bin
install:
	@echo "Installing binaries..."
	$(GO) install ./cmd/server
	$(GO) install ./cmd/client
	@echo "Installation complete!"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t open-source-vpn-server:latest -f docker/Dockerfile .
	@echo "Docker image built!"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Format complete!"

## lint: Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run ./...

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies updated!"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.DEFAULT_GOAL := help
