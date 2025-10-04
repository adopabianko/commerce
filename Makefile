.PHONY: up down proto-all build-all test-all

make up:
	docker compose up -d --build

make down:
	docker compose down -v
proto-all:
	cd proto && make

build-all:
	cd order-service && go build ./...
	cd inventory-service && go build ./...

test-all:
	cd order-service && go test ./...
	cd inventory-service && go test ./...
