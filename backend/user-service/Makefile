
.PHONY: build run migrate test clean docker-up init-db dev

build:
	go build -o bin/user-service.exe cmd/server/main.go

build-migrate:
	go build -o bin/migrate.exe cmd/migrate/main.go

docker-up:
	cd .. && docker-compose up -d

init-db: docker-up
	timeout /t 3
	docker exec -it postgres psql -U postgres -c "CREATE DATABASE bus_dev_user;" || echo "Database already exists"

dev: init-db migrate run

run:
	go run cmd/server/main.go

migrate: build-migrate
	bin\migrate.exe

test:
	go test -v ./...

clean:
	if exist bin rmdir /s /q bin

dev-setup: build-migrate migrate
	@echo "Development environment setup complete"

docker-build:
	docker build -t user-service .

docker-run:
	docker run -p 8080:8080 user-service

db-migrate: migrate

help:
	@echo "Available commands:"
	@echo "  build        - Build the main application"
	@echo "  build-migrate - Build the migration tool"
	@echo "  run          - Run the application"
	@echo "  migrate      - Run database migration"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  dev-setup    - Setup development environment"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  db-migrate   - Run database migration (alias for migrate)"
	@echo "  help         - Show this help message"