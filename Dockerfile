# Multi-stage build for jokes API
# Stage 1: Build the Go binary
<<<<<<< Updated upstream
FROM golang:1.21-alpine AS builder
=======
FROM golang:1.24-alpine AS builder
>>>>>>> Stashed changes

# Install build dependencies
RUN apk add --no-cache git make curl bash

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') \
    go build -ldflags="-s -w" -o jokes-api .

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Add labels
LABEL org.opencontainers.image.title="Jokes API"
LABEL org.opencontainers.image.description="5,160+ jokes across 16 categories"
LABEL org.opencontainers.image.vendor="APIMGR"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.url="https://jokes.apimgr.us"
LABEL org.opencontainers.image.source="https://github.com/apimgr/jokes"
LABEL org.opencontainers.image.version="1.0.0"

# Install runtime dependencies
RUN apk add --no-cache \
    curl \
    bash \
    ca-certificates \
    tzdata

# Create non-root user
RUN addgroup -g 1001 -S jokes && \
    adduser -S jokes -u 1001 -G jokes

# Create directories
RUN mkdir -p /data /config /var/log/jokes && \
    chown -R jokes:jokes /data /config /var/log/jokes

<<<<<<< Updated upstream
# Copy binary from builder
COPY --from=builder /build/jokes-api /usr/local/bin/jokes-api
RUN chmod +x /usr/local/bin/jokes-api

# Copy data files
COPY --chown=jokes:jokes src/data /data

=======
# Copy binary from builder (data is embedded in binary)
COPY --from=builder /build/jokes-api /usr/local/bin/jokes-api
RUN chmod +x /usr/local/bin/jokes-api

>>>>>>> Stashed changes
# Switch to non-root user
USER jokes

# Set working directory
WORKDIR /data

# Expose port 80 (internal)
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:80/healthz || exit 1

# Set environment variables
ENV GIN_MODE=release \
    PORT=80

# Run the application
ENTRYPOINT ["/usr/local/bin/jokes-api"]
CMD ["--port", "80", "--address", "0.0.0.0"]
