.PHONY: run build test clean fmt fmt-check vet lint quality-check swagger install-swag

run: ## Run the chatbot service
	go run cmd/server/main.go

build: ## Build the chatbot service
	go build -o bin/chatbot-service cmd/server/main.go

fmt: ## Format code
	go fmt ./...

fmt-check: ## Check if code is formatted
	@files=$$(gofmt -l .); \
	if [ -n "$$files" ]; then \
		echo "Code is not formatted. Run 'make fmt' to format."; \
		echo "$$files"; \
		exit 1; \
	fi

vet: ## Run go vet
	go vet ./...

lint: ## Run linter
	golangci-lint run --config=../.golangci.yml ./...

test: ## Run tests
	go test -v -cover ./...

quality-check: fmt-check vet lint test ## Run all quality checks

clean: ## Clean build files
	rm -rf bin/
	rm -rf docs/

# Swagger
install-swag: ## Install swag tool
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

swagger: ## Generate swagger documentation
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

swagger-gen: fmt swagger ## Format and generate swagger

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
