# Build stage
FROM golang:1.20.6-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o bookmyroom cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go cmd/web/send-mail.go
RUN go install github.com/gobuffalo/pop/v6/soda@latest

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/bookmyroom .
COPY --from=builder /go/bin/soda ./soda
COPY cmd/web/.env .
COPY migrations ./migrations
COPY database.yml .
COPY start.sh .
COPY wait-for.sh .

EXPOSE 9090
CMD ["/app/bookmyroom"]
ENTRYPOINT [ "/app/start.sh" ]