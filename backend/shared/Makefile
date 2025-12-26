.PHONY: run migrate fmt fmt-check fmt-fix lint test vet quality-check install-swag swagger swagger-gen

run:
	go run ./cmd/server/main.go

fmt:
	go fmt ./...

fmt-check:
	@files=$$(gofmt -l .); \
	if [ -n "$$files" ]; then \
		echo "Code is not formatted. Run 'make fmt' to format."; \
		echo "$$files"; \
		exit 1; \
	fi

lint:
	golangci-lint run --config=../.golangci.yml ./...

test:
	go test ./...

vet:
	go vet ./...

quality-check:
	fmt-check vet lint test

# Swagger
install-swag:
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

swagger-gen:
	fmt swagger

# Migration commands
DB_URL=postgres://doadmin:AVNS_k8dv8ULBQ3Sjwie2_Xe@bus-booking-database-do-user-29228381-0.d.db.ondigitalocean.com:25060/bus_booking_booking?sslmode=require

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-reset:
	migrate -path migrations -database "$(DB_URL)" drop -f
	migrate -path migrations -database "$(DB_URL)" up

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

# mock
# make mock-file dir=internal/repository name=bus_repository
mock-file:
	mockgen -source=$(dir)/$(name).go -package=mocks -destination=$(dir)/mocks/mock_$(name).go