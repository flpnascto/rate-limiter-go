FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o ratelimter ./cmd/main.go

FROM scratch
COPY --from=builder /app/ratelimter .
COPY --from=builder /app/cmd/.env .
CMD ["./ratelimter"]