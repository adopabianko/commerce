FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/user-service
RUN go mod download
RUN go build -o user-service ./cmd/server

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/user-service/user-service /user-service
EXPOSE 50051 9090
ENTRYPOINT ["/user-service"]
