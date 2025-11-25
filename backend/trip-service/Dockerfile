# Trip Service Dockerfile
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy only required modules (no go.work dependency)
COPY shared/ ./shared/
COPY trip-service/ ./trip-service/

WORKDIR /build/trip-service

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o app ./cmd/server

# Final stage
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/trip-service/app /app

ENV TZ=Asia/Ho_Chi_Minh

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["/app", "--health"] || exit 1

ENTRYPOINT ["/app"]