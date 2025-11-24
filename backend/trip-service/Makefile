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

.PHONY: build run migrate test clean docker-up init-db dev

build:
	go build -o bin/trip-service.exe cmd/server/main.go

build-migrate:
	go build -o bin/migrate.exe cmd/migrate/main.go

docker-up:
	cd .. && docker-compose up -d

init-db: docker-up
	timeout /t 3
	docker exec -it postgres psql -U postgres -c "CREATE DATABASE bus_dev_trip;" || echo "Database already exists"

dev: init-db migrate run


test:
	go test -v ./...

clean:
	if exist bin rmdir /s /q bin
	docker run -p 8081:8081 trip-service

# Database operations
db-reset: docker-up
	docker exec -it postgres psql -U postgres -c "DROP DATABASE IF EXISTS bus_dev_trip;"
	docker exec -it postgres psql -U postgres -c "CREATE DATABASE bus_dev_trip;"
	$(MAKE) migrate