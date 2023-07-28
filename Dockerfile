# Build stage
FROM golang:1.20.6-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o bookmyroom cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go cmd/web/send-mail.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/bookmyroom .

EXPOSE 8080
CMD ["/app/bookmyroom"]