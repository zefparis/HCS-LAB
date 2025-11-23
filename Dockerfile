# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hcsapi ./cmd/hcsapi

# Final stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/hcsapi .

# Expose port
EXPOSE 8080

# Set environment variable
ENV PORT=8080

# Run the application
CMD ["./hcsapi"]
