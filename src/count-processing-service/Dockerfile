# Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Enable GOTOOLCHAIN=auto to handle go version in go.mod
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/server/main.go

# Runner stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
