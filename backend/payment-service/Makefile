run:
	go run cmd/server/main.go

migrate:
	go run cmd/migrate/main.go

fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Code is not formatted. Run 'make fmt' to fix." && gofmt -l . && exit 1)

.PHONY: build run migrate test clean docker-up init-db dev lint fmt fmt-check vet quality-check

build:
	go build -o bin/payment-service.exe cmd/server/main.go

build-migrate:
	go build -o bin/migrate.exe cmd/migrate/main.go

docker-up:
	cd .. && docker-compose up -d

init-db: docker-up
	timeout /t 3
	docker exec -it postgres psql -U postgres -c "CREATE DATABASE bus_dev_payment;" || echo "Database already exists"

dev: init-db migrate run

test:
	go test -v ./...

clean:
	if exist bin rmdir /s /q bin

dev-setup: build-migrate migrate
	@echo "Development environment setup complete"

docker-build:
	docker build -t payment-service .

docker-run:
	docker run -p 8083:8083 payment-service

db-migrate: migrate

# Linting with golangci-lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config=../.golangci.yml ./...



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

help:
	@echo "Available commands:"
	@echo "  build        - Build the main application"
	@echo "  build-migrate - Build the migration tool"
	@echo "  run          - Run the application"
	@echo "  migrate      - Run database migration"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  lint         - Run golangci-lint"
	@echo "  fmt          - Fix code formatting"
	@echo "  fmt-check    - Check code formatting"
	@echo "  vet          - Run go vet"
	@echo "  quality-check - Run all quality checks"
	@echo "  dev-setup    - Setup development environment"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  db-migrate   - Run database migration (alias for migrate)"
	@echo "  help         - Show this help message"