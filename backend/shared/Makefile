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