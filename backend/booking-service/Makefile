run:
	go run cmd/server/main.go

migrate:
	go run cmd/migrate/main.go

fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Code is not formatted. Run 'make fmt' to fix." && gofmt -l . && exit 1)

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config=../.golangci.yml ./...

# Variables
BINARY_NAME=booking-service
MAIN_PATH=./cmd/server
MIGRATE_PATH=./cmd/migrate
BUILD_DIR=./bin
GO_FILES=$(shell find . -type f -name '*.go')

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-trimpath $(LDFLAGS)

.PHONY: all build build-linux clean test coverage fmt vet deps run dev migrate help

# Default target
all: clean deps fmt vet test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/migrate $(MIGRATE_PATH)
	@echo "Build completed successfully"

# Build for Linux
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/migrate-linux $(MIGRATE_PATH)
	@echo "Linux build completed successfully"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Tests completed"

# Run tests with coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format Go files
fmt:
	@echo "Formatting Go files..."
	@$(GOFMT) -s -w $(GO_FILES)
	@echo "Formatting completed"

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "Vet completed"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Run the application in development mode
dev: deps fmt vet
	@echo "Starting development server..."
	$(GOCMD) run $(MAIN_PATH)

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .
	@echo "Docker image built successfully"

docker-run:
	@echo "Running Docker container..."
	docker run -p 8082:8082 --env-file config/.env $(BINARY_NAME):latest

# Development with hot reload (requires air)
hot:
	@echo "Starting development server with hot reload..."
	@if ! command -v air > /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "Development tools installed"

# Static analysis
staticcheck:
	@echo "Running static analysis..."
	@if ! command -v staticcheck > /dev/null; then \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi
	staticcheck ./...
	@echo "Static analysis completed"

# Security scan
security:
	@echo "Running security scan..."
	@if ! command -v gosec > /dev/null; then \
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  coverage    - Run tests with coverage report"
	@echo "  fmt         - Format Go files"
	@echo "  vet         - Run go vet"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  dev         - Run in development mode"
	@echo "  run         - Build and run the application"
	@echo "  migrate     - Run database migrations"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  hot         - Start with hot reload (requires air)"
	@echo "  install-tools- Install development tools"
	@echo "  staticcheck - Run static analysis"
	@echo "  security    - Run security scan"
	@echo "  quality     - Run full quality check"
	@echo "  help        - Show this help message"