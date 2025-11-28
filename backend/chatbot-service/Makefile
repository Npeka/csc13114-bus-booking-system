.PHONY: help run build test clean lint fmt swagger migrate-up migrate-down

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run: ## Run the chatbot service
	go run cmd/server/main.go

build: ## Build the chatbot service
	go build -o bin/chatbot-service cmd/server/main.go

test: ## Run tests
	go test -v -cover ./...

clean: ## Clean build files
	rm -rf bin/
	rm -rf docs/

lint: ## Run linter
	golangci-lint run --config=../.golangci.yml ./...

fmt: ## Format code
	go fmt ./...

swagger: ## Generate swagger documentation
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
