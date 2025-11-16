# Bus Booking System - Infrastructure Management

.PHONY: help local-up local-down local-logs build-all deploy-k8s setup-argocd clean

# Default target
help: ## Show this help message
	@echo "Bus Booking System - Infrastructure Commands"
	@echo "=============================================="
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Local Development
local-up: ## Start all services locally with docker-compose
	@echo "🚀 Starting local development environment..."
	cd backend/infra && docker-compose up -d
	@echo "✅ Services started. Access at:"
	@echo "   - Gateway Service: http://localhost:8000 (Main Entry Point)"
	@echo "   - User Service: http://localhost:8080"
	@echo "   - Trip Service: http://localhost:8081" 
	@echo "   - Booking Service: http://localhost:8082"
	@echo "   - Template Service: http://localhost:8083"
	@echo "   - Payment Service: http://localhost:8084"
	@echo "   - PostgreSQL: localhost:5432"
	@echo "   - Redis: localhost:6379"

local-down: ## Stop all local services
	@echo "🛑 Stopping local development environment..."
	cd backend/infra && docker-compose down
	@echo "✅ All services stopped"

local-logs: ## Show logs for all local services
	cd backend/infra && docker-compose logs -f

local-restart: ## Restart specific service (usage: make local-restart SERVICE=user-service)
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please specify SERVICE name: make local-restart SERVICE=user-service"; \
		exit 1; \
	fi
	@echo "🔄 Restarting $(SERVICE)..."
	cd backend/infra && docker-compose restart $(SERVICE)

local-build: ## Rebuild and restart all services
	@echo "🔨 Building and restarting all services..."
	cd backend/infra && docker-compose up -d --build

# Docker Operations
build-all: ## Build all Docker images locally
	@echo "🔨 Building all Docker images..."
	@for service in user-service trip-service booking-service template-service payment-service gateway-service; do \
		echo "Building $$service..."; \
		cd backend && docker build -f $$service/Dockerfile -t bus-booking-$$service:latest .; \
	done
	@echo "✅ All images built successfully"

build-service: ## Build specific service (usage: make build-service SERVICE=user-service)
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please specify SERVICE name: make build-service SERVICE=user-service"; \
		exit 1; \
	fi
	@echo "🔨 Building $(SERVICE)..."
	cd backend && docker build -f $(SERVICE)/Dockerfile -t bus-booking-$(SERVICE):latest .
	@echo "✅ $(SERVICE) built successfully"

# Go Development
go-tidy: ## Run go mod tidy for all services
	@echo "📦 Running go mod tidy for all services..."
	@for service in booking-service trip-service user-service template-service payment-service gateway-service; do \
		echo "Tidying $$service..."; \
		cd backend/$$service && go mod tidy; \
	done
	@echo "✅ All go modules tidied"

go-test: ## Run tests for all services
	@echo "🧪 Running tests for all services..."
	@for service in booking-service trip-service user-service template-service payment-service gateway-service; do \
		echo "Testing $$service..."; \
		cd backend/$$service && go test ./...; \
	done

go-build: ## Build all Go binaries
	@echo "🔨 Building all Go binaries..."
	@for service in booking-service trip-service user-service template-service payment-service gateway-service; do \
		echo "Building $$service..."; \
		cd backend/$$service && go build -o bin/server ./cmd/server; \
	done
	@echo "✅ All binaries built"

# Kubernetes Operations  
k8s-deploy: ## Deploy to Kubernetes using Helm
	@echo "🚀 Deploying to Kubernetes..."
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm upgrade --install bus-booking backend/infra/k8s/helm/bus-booking \
		--namespace bus-booking \
		--create-namespace \
		--values backend/infra/k8s/helm/bus-booking/values.yaml
	@echo "✅ Deployed to Kubernetes"

k8s-status: ## Check Kubernetes deployment status
	@echo "📊 Checking Kubernetes status..."
	kubectl get pods -n bus-booking
	kubectl get services -n bus-booking
	kubectl get ingress -n bus-booking

k8s-logs: ## Get logs from Kubernetes pods
	@echo "📋 Getting pod logs..."
	kubectl logs -l app.kubernetes.io/name=bus-booking -n bus-booking --tail=100

k8s-delete: ## Delete Kubernetes deployment
	@echo "🗑️ Deleting Kubernetes deployment..."
	helm uninstall bus-booking -n bus-booking
	kubectl delete namespace bus-booking

# ArgoCD Operations
argocd-setup: ## Setup ArgoCD application
	@echo "🔧 Setting up ArgoCD application..."
	kubectl apply -f backend/infra/k8s/argocd/application.yaml
	@echo "✅ ArgoCD application created"

argocd-sync: ## Sync ArgoCD application
	@echo "🔄 Syncing ArgoCD application..."
	argocd app sync bus-booking-system
	@echo "✅ ArgoCD sync completed"

argocd-status: ## Check ArgoCD application status
	argocd app get bus-booking-system

# Utility Commands
clean: ## Clean up Docker resources
	@echo "🧹 Cleaning up Docker resources..."
	docker system prune -f
	docker volume prune -f
	@echo "✅ Cleanup completed"

setup-secrets: ## Setup GitHub secrets (requires gh CLI)
	@echo "🔐 Setting up GitHub secrets..."
	@echo "Please run these commands manually:"
	@echo "gh secret set DOCKER_USERNAME --body 'your-dockerhub-username'"
	@echo "gh secret set DOCKER_PASSWORD --body 'your-dockerhub-token'"
	@echo "gh secret set ARGOCD_SERVER --body 'https://your-argocd-server'"
	@echo "gh secret set ARGOCD_TOKEN --body 'your-argocd-token'"

health-check: ## Check health of all services
	@echo "🏥 Checking service health..."
	@for port in 8000 8080 8081 8082 8083 8084; do \
		echo -n "Port $$port: "; \
		curl -s -o /dev/null -w "%{http_code}" http://localhost:$$port/health || echo "❌ Not responding"; \
		echo ""; \
	done

# Database Operations
db-migrate: ## Run database migrations for all services
	@echo "🗃️ Running database migrations..."
	@for service in booking-service trip-service user-service; do \
		echo "Migrating $$service..."; \
		cd backend/$$service && go run cmd/migrate/main.go; \
	done
	@echo "✅ Migrations completed"

db-seed: ## Seed database with test data
	@echo "🌱 Seeding database..."
	cd backend/infra && docker-compose exec postgres psql -U postgres -d postgres -f /docker-entrypoint-initdb.d/init-db.sql
	@echo "✅ Database seeded"

# Development Helpers
dev-setup: ## Setup development environment
	@echo "🔧 Setting up development environment..."
	@echo "1. Installing dependencies..."
	$(MAKE) go-tidy
	@echo "2. Building services..."
	$(MAKE) go-build
	@echo "3. Starting services..."
	$(MAKE) local-up
	@echo "✅ Development environment ready!"

# Production Helpers
prod-deploy: ## Full production deployment
	@echo "🚀 Starting production deployment..."
	@echo "1. Building images..."
	$(MAKE) build-all
	@echo "2. Deploying to Kubernetes..."
	$(MAKE) k8s-deploy
	@echo "3. Setting up ArgoCD..."
	$(MAKE) argocd-setup
	@echo "✅ Production deployment completed!"

# CI/CD Helpers
ci-test: ## Run CI tests locally
	@echo "🧪 Running CI tests..."
	$(MAKE) go-tidy
	$(MAKE) go-test
	$(MAKE) build-all
	@echo "✅ CI tests completed"

# Documentation
docs: ## Generate API documentation
	@echo "📚 Generating API documentation..."
	@echo "TODO: Add swagger generation"

# Version Management
version: ## Show current version
	@echo "Bus Booking System v1.0.0"
	@echo "Components:"
	@echo "  - Go: $(shell go version)"
	@echo "  - Docker: $(shell docker --version)"
	@echo "  - Kubernetes: $(shell kubectl version --client --short 2>/dev/null || echo 'Not installed')"
	@echo "  - Helm: $(shell helm version --short 2>/dev/null || echo 'Not installed')"