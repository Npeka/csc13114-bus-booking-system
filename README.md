# Bus Booking System

A comprehensive microservices-based bus booking platform built with modern technologies, featuring a Next.js frontend and Go-based backend services.

## 🏗️ Architecture Overview

This project implements a distributed microservices architecture with the following components:

## ✅ Week 3 Completion Status

**Booking System and Seat Selection - Completed**

The following features have been successfully implemented as part of Week 3:

### User Portal Features

- ✅ **Interactive Seat Selection**

  - Visual seat map component with clickable seats
  - Different seat types and status indicators
  - Real-time availability updates
  - Seat locking mechanism to prevent double bookings

- ✅ **Booking Flow**

  - Passenger information collection forms
  - Booking creation and management APIs
  - Booking summary and review interface
  - Booking history and management dashboard
  - Booking modification and cancellation

- ✅ **Guest Services**

  - Guest checkout without registration
  - Guest booking lookup by reference/email/phone
  - Unique booking reference generation

- ✅ **Ticketing**
  - E-ticket download functionality
  - Automatic email delivery of e-tickets
  - Professional e-ticket template with branding

### Frontend Implementation

- Profile management page with view/edit modes
- Real-time trip search with filters
- Booking management dashboard
- Responsive UI components

### Backend APIs

- Complete booking service with PostgreSQL
- Seat reservation and locking system
- Trip and route management
- User profile management
- Gateway authentication and routing

**Next Up:** Week 4 focuses on payment integration (PayOS), notifications (email/SMS), and admin analytics.

See [NEXT_STEPS.md](./NEXT_STEPS.md) for detailed Week 4 planning.

### Backend Services (Go)

- **User Service** - Authentication, user management, and Firebase integration ✅ **Fully Implemented (Week 1-2)**
- **Trip Service** - Bus trip management and scheduling ✅ **Fully Implemented (Week 2)**
- **Booking Service** - Reservation and booking management ✅ **Fully Implemented (Week 3)**
- **Payment Service** - Payment processing and transaction handling 🚧 **Week 4 - In Progress**
- **Gateway Service** - API Gateway for routing and request aggregation ✅ **Operational**

> **Note:** The project is currently at **Week 3 completion**. User Service, Trip Service, and Booking Service are fully implemented. Payment integration is planned for Week 4.

### Frontend (Next.js)

- Modern React-based web application with TypeScript
- Server-side rendering and static generation
- Responsive UI with Tailwind CSS and Radix UI components

### Infrastructure

- **PostgreSQL** - Single PostgreSQL server with one database per service
- **Redis** - Caching and session management
- **Docker** - Containerization for all services
- **DigitalOcean Kubernetes** - Production deployment orchestration
- **ArgoCD** - GitOps-based continuous deployment using Helm charts from `infra` branch
- **Ingress + Cert Manager** - Traffic routing and SSL certificate management

## 🌐 API Gateway Architecture

All external traffic flows through the **Gateway Service** which acts as a single entry point:

1. **Ingress** routes all traffic to the Gateway Service
2. **Gateway Service** receives requests and checks route configuration
3. For protected routes, Gateway calls **User Service** to verify JWT token
4. If authentication succeeds, Gateway forwards the request to the appropriate backend service
5. User context (ID, email, role) is added to request headers for downstream services

This architecture ensures:

- Centralized authentication and authorization
- Service-to-service communication remains internal
- Role-based access control (RBAC) enforcement
- Clean separation of concerns

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
│   ├── gateway-service/       # API Gateway with JWT verification
│   ├── payment-service/       # Payment processing
│   ├── trip-service/          # Trip management
│   ├── user-service/          # User authentication & management
│   ├── shared/                # Shared libraries and utilities
│   ├── init-databases.sql     # PostgreSQL database initialization
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

The project uses **Husky** for Git hooks and **lint-staged** for running checks on staged files.

#### Setup

Pre-commit hooks are automatically installed when you run:

```bash
npm install
```

This runs the `prepare` script which executes `husky install`.

#### What Gets Checked

**Frontend** (when `frontend/` files are staged):

- ESLint checks via `pnpm lint:check`
- Prettier formatting via `pnpm format:check`
- Runs through lint-staged for efficiency

**Backend** (per-service, only for changed services):

- `gofmt` formatting check
- `golangci-lint` static analysis
- Automatically detects which services changed
- Only runs checks for modified services

#### Example Output

```bash
[PRE] Pre-commit started
[FE ] frontend changed → running lint-staged...
✔ Running tasks for staged files...
[BE ] Running Go fmt + lint for changed backend services...
[BE ] booking-service no changes → skip
[BE ] payment-service changed → running fmt + lint...
golangci-lint run --config=../.golangci.yml ./...
0 issues.
[BE ] user-service no changes → skip
[OK ] Pre-commit checks passed.
```

#### Configuration Files

- [`.husky/pre-commit`](.husky/pre-commit) - Main pre-commit hook script
- [`package.json`](package.json) - Husky and lint-staged configuration
- Root level manages both frontend and backend checks

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

## ☸️ Production Deployment

### Deployment Architecture

The project uses a GitOps approach with the following components:

- **Platform**: DigitalOcean Kubernetes
- **GitOps Tool**: ArgoCD
- **Package Manager**: Helm
- **Infrastructure Config**: Stored in `infra` branch (separate from `main`)
- **Ingress Controller**: Routes external traffic to Gateway Service only
- **Cert Manager**: Automatic SSL certificate management

### Deployment Flow

1. **Code Push**: Developer pushes code to `main` or `develop` branch
2. **Quality Checks**: GitHub Actions runs linting, formatting, and tests
3. **Build & Push**: Docker images are built and pushed to Docker Hub with versioned tags
4. **Update Infra**: GitHub Actions updates Helm values in `infra` branch with new image tags
5. **ArgoCD Sync**: ArgoCD detects changes in `infra` branch and deploys automatically
6. **Kubernetes**: Pods are updated with new images via rolling update

### CI/CD Pipeline

See [`.github/workflows/quality-checks.yml`](.github/workflows/quality-checks.yml) and [`.github/workflows/build-deploy.yml`](.github/workflows/build-deploy.yml) for the complete pipeline configuration.

## 🗄️ Database

### Database Architecture

- **Single PostgreSQL server** with one database per microservice
- Separate database per service for data isolation and independent scaling
- Automatic initialization via [`init-databases.sql`](backend/init-databases.sql)

### Databases

- `bus_booking_user` - User service database
- `bus_booking_trip` - Trip service database
- `bus_booking_booking` - Booking service database
- `bus_booking_payment` - Payment service database

### Database Initialization

Databases are automatically created when the PostgreSQL container starts for the first time using the [`init-databases.sql`](backend/init-databases.sql) script.

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
