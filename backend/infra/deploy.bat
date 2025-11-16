@echo off
setlocal enabledelayedexpansion

rem Quick deployment script for Bus Booking System (Windows)
rem Usage: deploy.bat [local|k8s|clean]

set "SCRIPT_DIR=%~dp0"
set "PROJECT_ROOT=%SCRIPT_DIR%..\.."

rem Default to local deployment
set "COMMAND=%~1"
if "%COMMAND%"=="" set "COMMAND=local"

rem Colors (if supported)
set "INFO=[INFO]"
set "SUCCESS=[SUCCESS]"
set "WARNING=[WARNING]"
set "ERROR=[ERROR]"

echo %INFO% Bus Booking System Deployment Script

rem Check prerequisites
echo %INFO% Checking prerequisites...
docker --version >nul 2>&1
if errorlevel 1 (
    echo %ERROR% Docker is not installed
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo %ERROR% Docker Compose is not installed
    exit /b 1
)

echo %SUCCESS% Prerequisites check passed

rem Execute command
if /i "%COMMAND%"=="local" goto :deploy_local
if /i "%COMMAND%"=="k8s" goto :deploy_k8s
if /i "%COMMAND%"=="kubernetes" goto :deploy_k8s
if /i "%COMMAND%"=="build" goto :build_all
if /i "%COMMAND%"=="clean" goto :clean_up

echo Usage: %0 [local^|k8s^|build^|clean]
echo.
echo Commands:
echo   local   - Deploy locally with Docker Compose (default)
echo   k8s     - Deploy to Kubernetes with Helm
echo   build   - Build all Docker images
echo   clean   - Clean up local resources
exit /b 1

:deploy_local
echo %INFO% Starting local deployment...
cd /d "%SCRIPT_DIR%"

rem Start services
docker-compose up -d
if errorlevel 1 (
    echo %ERROR% Failed to start services
    exit /b 1
)

rem Wait for services
echo %INFO% Waiting for services to be healthy...
timeout /t 30 /nobreak >nul

rem Check service health
echo %INFO% Checking service health...
for %%p in (8080 8081 8082 8083 8084) do (
    curl -f -s http://localhost:%%p/health >nul 2>&1
    if errorlevel 1 (
        echo %WARNING% Service on port %%p is not responding
    ) else (
        echo %SUCCESS% Service on port %%p is healthy
    )
)

echo %SUCCESS% Local deployment completed!
echo %INFO% Services are available at:
echo   - User Service: http://localhost:8080
echo   - Trip Service: http://localhost:8081
echo   - Booking Service: http://localhost:8082
echo   - Template Service: http://localhost:8083
echo   - Payment Service: http://localhost:8084
echo   - PostgreSQL: localhost:5432
echo   - Redis: localhost:6379
goto :end

:deploy_k8s
echo %INFO% Starting Kubernetes deployment...

rem Check kubectl
kubectl version --client >nul 2>&1
if errorlevel 1 (
    echo %ERROR% kubectl is not installed
    exit /b 1
)

rem Check helm
helm version >nul 2>&1
if errorlevel 1 (
    echo %ERROR% Helm is not installed
    exit /b 1
)

rem Add Bitnami repository
echo %INFO% Adding Bitnami Helm repository...
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

rem Deploy with Helm
echo %INFO% Deploying with Helm...
helm upgrade --install bus-booking "%SCRIPT_DIR%k8s\helm\bus-booking" ^
    --namespace bus-booking ^
    --create-namespace ^
    --values "%SCRIPT_DIR%k8s\helm\bus-booking\values.yaml"

if errorlevel 1 (
    echo %ERROR% Helm deployment failed
    exit /b 1
)

rem Wait for deployment
echo %INFO% Waiting for deployment to be ready...
kubectl wait --for=condition=available --timeout=300s deployment -l app.kubernetes.io/name=bus-booking -n bus-booking

rem Show status
echo %INFO% Deployment status:
kubectl get pods -n bus-booking
kubectl get services -n bus-booking

echo %SUCCESS% Kubernetes deployment completed!
goto :end

:build_all
echo %INFO% Building all services...
cd /d "%PROJECT_ROOT%"

for %%s in (user-service trip-service booking-service template-service payment-service) do (
    echo %INFO% Building %%s...
    cd backend
    docker build -f "%%s\Dockerfile" -t "bus-booking-%%s:latest" .
    if errorlevel 1 (
        echo %ERROR% Failed to build %%s
        exit /b 1
    )
    cd ..
)

echo %SUCCESS% All services built successfully!
goto :end

:clean_up
echo %INFO% Cleaning up...
cd /d "%SCRIPT_DIR%"

rem Stop local services
docker-compose down -v

rem Clean Docker resources
docker system prune -f
docker volume prune -f

echo %SUCCESS% Cleanup completed!
goto :end

:end
echo %INFO% Script completed
exit /b 0