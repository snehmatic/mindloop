# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build CLI and Server binaries
RUN CGO_ENABLED=1 GOOS=linux go build -o mindloop main.go
RUN CGO_ENABLED=1 GOOS=linux go build -o mindloop-server cmd/server/server.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache libc6-compat ca-certificates

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/mindloop /usr/local/bin/mindloop
COPY --from=builder /app/mindloop-server /usr/local/bin/mindloop-server

# Copy web assets
COPY --from=builder /app/web ./web

# Expose the default port
EXPOSE 8765

# Set the entrypoint to the server binary
ENTRYPOINT ["mindloop-server"]

# Default arguments for the server
CMD ["-port=8765", "-mode=local"]
