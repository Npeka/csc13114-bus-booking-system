.PHONY: run migrate test lint fmt fmt-fix vet quality-check

run:
	go run cmd/server/main.go

migrate:
	go run cmd/migrate/main.go

test:
	go test -v ./...

# Linting with golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config=../.golangci.yml ./...

# Check formatting
fmt:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Code is not formatted. Run 'make fmt-fix' to fix." && gofmt -l . && exit 1)

# Fix formatting
fmt-fix:
	@echo "Fixing code formatting..."
	@gofmt -w .
.PHONY: run migrate test lint fmt fmt-fix vet quality-check install-swag swagger swagger-gen

run:
	go run cmd/server/main.go

migrate:
	go run cmd/migrate/main.go

test:
	go test -v ./...

# Linting with golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config=../.golangci.yml ./...

# Check formatting
fmt:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Code is not formatted. Run 'make fmt-fix' to fix." && gofmt -l . && exit 1)

# Fix formatting
fmt-fix:
	@echo "Fixing code formatting..."
	@gofmt -w .

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run all quality checks
quality-check: fmt vet lint test
	@echo "All quality checks passed!"

# Install swag tool for Swagger documentation generation
install-swag:
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Swag installed successfully!"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger documentation generated in docs/ directory"

# Format and generate Swagger docs
swagger-gen: fmt swagger
	@echo "Swagger documentation ready!"