# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o health-checker .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/health-checker .

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080 \
    DATABASE_URL=http://database:5432 \
    CACHE_URL=http://redis:6379 \
    API_URL=http://api:8000 \
    MONITORING_URL=http://prometheus:9090

# Run the application
CMD ["./health-checker"]
