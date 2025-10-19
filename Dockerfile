# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o show-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/show-service .

# Copy config files
COPY --from=builder /app/configs ./configs

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./show-service"]
