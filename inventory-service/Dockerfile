FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/inventory-service
RUN go mod download
RUN go build -o inventory-service ./cmd/server

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/inventory-service/inventory-service /inventory-service
EXPOSE 50051 9090
ENTRYPOINT ["/inventory-service"]
