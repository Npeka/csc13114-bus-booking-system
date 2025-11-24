.PHONY: run test lint fmt fmt-check vet quality-check

run:
	go run cmd/server/main.go

test:
	go test -v ./...

# Linting with golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config=../.golangci.yml ./...

# Check formatting
fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Code is not formatted. Run 'make fmt' to fix." && gofmt -l . && exit 1)

# Fix formatting
fmt:
	@echo "Fixing code formatting..."
	@gofmt -w .

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run all quality checks
quality-check: fmt-check vet lint test
	@echo "All quality checks passed!"