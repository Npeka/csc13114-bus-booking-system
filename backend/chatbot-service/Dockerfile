# Chatbot Service Dockerfile
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata wget
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy only required modules (matching payment-service pattern)
COPY shared/ ./shared/
COPY chatbot-service/ ./chatbot-service/

WORKDIR /build/chatbot-service

# Download dependencies
RUN go mod download

# Build the binary with optimized flags
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o app ./cmd/server

# Final stage - use scratch for minimal image size
FROM alpine:latest

# Install wget for healthcheck
RUN apk add --no-cache ca-certificates wget

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/chatbot-service/app /app

ENV TZ=Asia/Ho_Chi_Minh

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app"]
