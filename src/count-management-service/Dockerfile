# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o count-management-service .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and templates
COPY --from=builder /app/count-management-service .
COPY --from=builder /app/templates ./templates

# Expose port
EXPOSE 8080

# Run the application
CMD ["./count-management-service"]
