# Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Set GOTOOLCHAIN=auto as per instructions
ENV GOTOOLCHAIN=auto

# Download dependencies
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o count-api-service ./cmd/count-api-service/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/count-api-service .

# Create data directory
RUN mkdir -p /data
ENV STORAGE_PATH=/data/counts.log
ENV PORT=8080

EXPOSE 8080

CMD ["./count-api-service"]
