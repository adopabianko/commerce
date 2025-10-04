FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/order-service
RUN go mod download
RUN go build -o order-service ./cmd/server

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/order-service/order-service /order-service
EXPOSE 8080
ENTRYPOINT ["/order-service"]
