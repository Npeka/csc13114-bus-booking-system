.PHONY: run migrate fmt fmt-check fmt-fix lint test vet quality-check install-swag swagger swagger-gen

run:
	go run ./cmd/server/main.go

migrate:
	go run ./cmd/migrate/main.go

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
SWAG_PKG=github.com/swaggo/swag/v2/cmd/swag
SWAG_VERSION=v2.6.2

install-swag:
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

swagger-gen:
	fmt swagger