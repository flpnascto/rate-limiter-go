FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ratelimiter ./cmd/main.go