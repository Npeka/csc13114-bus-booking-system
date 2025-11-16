# User Service Dockerfile
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy workspace files
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY user-service/ ./user-service/

# Copy go.mod files
COPY user-service/go.mod user-service/go.sum ./user-service/
COPY shared/go.mod shared/go.sum ./shared/

WORKDIR /build/user-service

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o app ./cmd/server

# Final stage
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/user-service/app /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["/app", "--health"] || exit 1

ENTRYPOINT ["/app"]