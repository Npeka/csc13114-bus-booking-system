#!/bin/bash

# Quick deployment script for Bus Booking System
# Usage: ./deploy.sh [local|k8s|clean]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Local deployment
deploy_local() {
    log_info "Starting local deployment..."
    
    cd "$SCRIPT_DIR"
    
    # Start services
    docker-compose up -d
    
    # Wait for services to be healthy
    log_info "Waiting for services to be healthy..."
    sleep 30
    
    # Check service health
    log_info "Checking service health..."
    for port in 8080 8081 8082 8083 8084; do
        if curl -f -s http://localhost:$port/health > /dev/null; then
            log_success "Service on port $port is healthy"
        else
            log_warning "Service on port $port is not responding"
        fi
    done
    
    log_success "Local deployment completed!"
    log_info "Services are available at:"
    echo "  - User Service: http://localhost:8080"
    echo "  - Trip Service: http://localhost:8081"
    echo "  - Booking Service: http://localhost:8082"
    echo "  - Template Service: http://localhost:8083"
    echo "  - Payment Service: http://localhost:8084"
    echo "  - PostgreSQL: localhost:5432"
    echo "  - Redis: localhost:6379"
}

# Kubernetes deployment
deploy_k8s() {
    log_info "Starting Kubernetes deployment..."
    
    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed"
        exit 1
    fi
    
    # Check helm
    if ! command -v helm &> /dev/null; then
        log_error "Helm is not installed"
        exit 1
    fi
    
    # Add Bitnami repository
    log_info "Adding Bitnami Helm repository..."
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo update
    
    # Deploy with Helm
    log_info "Deploying with Helm..."
    helm upgrade --install bus-booking "$SCRIPT_DIR/k8s/helm/bus-booking" \
        --namespace bus-booking \
        --create-namespace \
        --values "$SCRIPT_DIR/k8s/helm/bus-booking/values.yaml"
    
    # Wait for deployment
    log_info "Waiting for deployment to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment -l app.kubernetes.io/name=bus-booking -n bus-booking
    
    # Show status
    log_info "Deployment status:"
    kubectl get pods -n bus-booking
    kubectl get services -n bus-booking
    
    log_success "Kubernetes deployment completed!"
}

# Clean up
clean_up() {
    log_info "Cleaning up..."
    
    # Stop local services
    cd "$SCRIPT_DIR"
    docker-compose down -v
    
    # Clean Docker resources
    docker system prune -f
    docker volume prune -f
    
    log_success "Cleanup completed!"
}

# Build all services
build_all() {
    log_info "Building all services..."
    
    cd "$PROJECT_ROOT"
    
    for service in user-service trip-service booking-service template-service payment-service; do
        log_info "Building $service..."
        cd "backend"
        docker build -f "$service/Dockerfile" -t "bus-booking-$service:latest" .
        cd ..
    done
    
    log_success "All services built successfully!"
}

# Main function
main() {
    case "${1:-local}" in
        local)
            check_prerequisites
            deploy_local
            ;;
        k8s|kubernetes)
            check_prerequisites
            deploy_k8s
            ;;
        build)
            check_prerequisites
            build_all
            ;;
        clean)
            clean_up
            ;;
        *)
            echo "Usage: $0 [local|k8s|build|clean]"
            echo ""
            echo "Commands:"
            echo "  local   - Deploy locally with Docker Compose (default)"
            echo "  k8s     - Deploy to Kubernetes with Helm"
            echo "  build   - Build all Docker images"
            echo "  clean   - Clean up local resources"
            exit 1
            ;;
    esac
}

main "$@"