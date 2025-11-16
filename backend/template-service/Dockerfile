# Multi-stage build for template-service
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates for go modules
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go.work and shared module first for better layer caching
COPY go.work go.work.sum ./
COPY shared/ ./shared/

# Copy service-specific files
COPY template-service/go.mod template-service/go.sum ./template-service/
WORKDIR /app/template-service

# Download dependencies
RUN go mod download

# Copy service source code
COPY template-service/ .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o server ./cmd/server

# Final stage - minimal image
FROM scratch

# Import ca-certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Create non-root user structure
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the binary
COPY --from=builder /app/template-service/server /server

# Set non-root user
USER 65534:65534

# Expose port
ARG SERVICE_PORT=8083
EXPOSE ${SERVICE_PORT}

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/server", "--health-check"]

ENTRYPOINT ["/server"]