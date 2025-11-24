# Bus Booking System

A comprehensive microservices-based bus booking platform built with modern technologies, featuring a Next.js frontend and Go-based backend services.

## 🏗️ Architecture Overview

This project implements a distributed microservices architecture with the following components:

### Backend Services (Go)
- **User Service** - Authentication, user management, and Firebase integration ✅ **Fully Implemented**
- **Trip Service** - Bus trip management and scheduling ⚠️ *Base project only*
- **Booking Service** - Reservation and booking management ⚠️ *Base project only*
- **Payment Service** - Payment processing and transaction handling ⚠️ *Base project only*
- **Gateway Service** - API Gateway for routing and request aggregation ⚠️ *Base project only*

> **Note:** Currently, only the **User Service** is fully implemented according to the project requirements. Other backend services have been scaffolded with base project structure for future development.

### Frontend (Next.js)
- Modern React-based web application with TypeScript
- Server-side rendering and static generation
- Responsive UI with Tailwind CSS and Radix UI components

### Infrastructure
- **PostgreSQL** - Primary database with separate schemas per service
- **Redis** - Caching and session management
- **Docker** - Containerization for all services
- **Kubernetes** - Production deployment orchestration
- **ArgoCD** - GitOps-based continuous deployment

## 🚀 Quick Start

### Prerequisites
- **Go** 1.24.0 or higher
- **Node.js** 20.x or higher
- **pnpm** 10.x
- **Docker** and **Docker Compose**
- **Make** (optional, for convenience commands)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd csc13114-bus-booking-system
   ```

2. **Start all services with Docker Compose**
   ```bash
   make local-up
   # or manually:
   cd backend && docker-compose up -d
   ```

3. **Access the services**
   - Gateway Service: http://localhost:8000 (Main Entry Point)
   - User Service: http://localhost:8080
   - Trip Service: http://localhost:8081
   - Booking Service: http://localhost:8082
   - Payment Service: http://localhost:8084
   - PostgreSQL: localhost:5432
   - Redis: localhost:6379

4. **Run the frontend**
   ```bash
   cd frontend
   pnpm install
   pnpm dev
   ```

### Stop Services
```bash
make local-down
# or manually:
cd backend && docker-compose down
```

## 📁 Project Structure

```
csc13114-bus-booking-system/
├── backend/
│   ├── booking-service/       # Booking management microservice
│   ├── gateway-service/       # API Gateway
│   ├── payment-service/       # Payment processing
│   ├── trip-service/          # Trip management
│   ├── user-service/          # User authentication & management
│   ├── shared/                # Shared libraries and utilities
│   ├── infra/                 # Infrastructure configurations
│   │   └── docker/            # Docker and database initialization
│   ├── docker-compose.yaml    # Local development orchestration
│   └── .golangci.yml          # Go linting configuration
├── frontend/
│   ├── app/                   # Next.js app directory
│   ├── components/            # React components
│   ├── lib/                   # Utility libraries
│   └── __tests__/             # Frontend tests
├── .github/
│   └── workflows/             # CI/CD pipelines
│       ├── quality-checks.yml # Linting, formatting, and tests
│       └── build-deploy.yml   # Build and deployment pipeline
└── Makefile                   # Convenience commands
```

## 🔧 Development

### Backend Development

Each Go microservice follows a clean architecture pattern with:
- **cmd/** - Application entry points
- **internal/** - Private application code
  - **handlers/** - HTTP request handlers
  - **services/** - Business logic
  - **repositories/** - Data access layer
  - **models/** - Data models
  - **utils/** - Utility functions
- **docs/** - Swagger API documentation

#### Common Commands

```bash
# Run all tests
make go-test

# Build all services
make go-build

# Tidy dependencies
make go-tidy

# Build specific service
make build-service SERVICE=user-service

# Restart specific service
make local-restart SERVICE=user-service
```

### Frontend Development

The frontend is built with:
- **Next.js 16** with App Router
- **TypeScript** for type safety
- **Tailwind CSS** for styling
- **Radix UI** for accessible components
- **React Query** for data fetching
- **Zustand** for state management
- **React Hook Form** with Zod validation

#### Available Scripts

```bash
cd frontend

# Development server
pnpm dev

# Build for production
pnpm build

# Run tests
pnpm test

# Lint code
pnpm lint

# Format code
pnpm format
```

## 🔍 Quality Assurance

### Automated Quality Checks

The project uses GitHub Actions for automated quality checks on every push and pull request:

#### Frontend Quality Checks
- **Type checking** - TypeScript compilation
- **Linting** - ESLint with Next.js configuration
- **Formatting** - Prettier code formatting
- **Tests** - Jest unit tests

#### Backend Quality Checks
Per-service checks with change detection:
- **Format checking** - `gofmt` compliance
- **Go vet** - Static analysis
- **Linting** - `golangci-lint` with custom rules
- **Unit tests** - Go test suite with coverage

### Pre-commit Hooks

The project uses Husky and lint-staged for pre-commit validation:
- Automatic code formatting
- Linting on staged files
- Type checking

### CI/CD Pipeline

The deployment pipeline consists of two workflows:

1. **Quality Checks** (`.github/workflows/quality-checks.yml`)
   - Runs on every push to `main` or `develop`
   - Detects changed services
   - Runs quality checks only for affected services
   - Must pass before build and deploy

2. **Build and Deploy** (`.github/workflows/build-deploy.yml`)
   - Triggers only after quality checks pass
   - Builds Docker images for changed services
   - Pushes to Docker Hub with versioned tags
   - Updates infrastructure repository (infra branch)
   - ArgoCD automatically deploys changes

## 🐳 Docker

### Build All Images
```bash
make build-all
```

### Build Specific Service
```bash
make build-service SERVICE=user-service
```

### Docker Compose Features
- Health checks for all services
- Automatic database initialization
- Service dependency management
- Network isolation
- Volume persistence for databases

## ☸️ Kubernetes Deployment

### Deploy to Kubernetes
```bash
make k8s-deploy
```

### Check Deployment Status
```bash
make k8s-status
```

### View Logs
```bash
make k8s-logs
```

### ArgoCD Integration
```bash
# Setup ArgoCD application
make argocd-setup

# Sync application
make argocd-sync

# Check status
make argocd-status
```

## 🗄️ Database

### Database Architecture
- Single PostgreSQL instance with multiple databases
- Separate database per microservice for data isolation
- Automatic initialization via init scripts

### Databases
- `user_db` - User service database
- `trip_db` - Trip service database
- `booking_db` - Booking service database
- `payment_db` - Payment service database

### Database Operations
```bash
# Run migrations
make db-migrate

# Seed test data
make db-seed
```

## 🔐 Environment Variables

Each service requires environment configuration. Example `.env.dev` files are provided in each service directory.

### Required Variables
- Database connection strings
- Redis connection
- JWT secrets
- Firebase credentials
- Service URLs

## 📊 Monitoring & Health Checks

All services expose health check endpoints:
```bash
# Check all service health
make health-check
```

Health endpoints:
- Gateway: `http://localhost:8000/health`
- User Service: `http://localhost:8080/health`
- Trip Service: `http://localhost:8081/health`
- Booking Service: `http://localhost:8082/health`
- Payment Service: `http://localhost:8084/health`

## 📚 API Documentation

Each backend service includes Swagger documentation:
- User Service: `http://localhost:8080/swagger/index.html`
- Trip Service: `http://localhost:8081/swagger/index.html`
- Booking Service: `http://localhost:8082/swagger/index.html`
- Payment Service: `http://localhost:8084/swagger/index.html`

## 🧪 Testing

### Backend Tests
```bash
# Run all service tests
make go-test

# Run tests for specific service
cd backend/user-service && go test ./...
```

### Frontend Tests
```bash
cd frontend

# Run tests
pnpm test

# Run tests in watch mode
pnpm test:watch

# Generate coverage report
pnpm test:coverage
```

## 🛠️ Troubleshooting

### Services won't start
```bash
# Clean up Docker resources
make clean

# Rebuild all services
make local-build
```

### Database connection issues
- Ensure PostgreSQL is healthy: `docker ps`
- Check database initialization logs: `docker logs postgres`
- Verify environment variables in `.env.dev` files

### Port conflicts
- Check if ports are already in use
- Modify port mappings in `docker-compose.yaml`

## 📝 Contributing

1. Create a feature branch from `develop`
2. Make your changes
3. Ensure all quality checks pass locally
4. Submit a pull request to `develop`
5. Wait for CI/CD pipeline to complete
6. Request code review

## 📄 License

This project is part of the CSC13114 course.

## 👥 Team

Developed as part of the CSC13114 - Advanced Web Development course.

---

For detailed documentation on individual services, see:
- [Frontend README](./frontend/README.md)
- [Backend Service Documentation](./backend/)
