# Build stage
FROM golang:1.27-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build CLI and Server binary
RUN CGO_ENABLED=1 GOOS=linux go build -o mindloop main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache libc6-compat ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/mindloop /usr/local/bin/mindloop

# Expose the default port
EXPOSE 8765

# Set the entrypoint to the mindloop binary
ENTRYPOINT ["mindloop"]

# Default arguments to start the server
CMD ["server", "--port=8765"]