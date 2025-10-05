PROTOC ?= protoc
OUT := ./gen

all: clean gen

gen: gen-inventory gen-user

gen-inventory:
	mkdir -p $(OUT)
	$(PROTOC) \
		--go_out=$(OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT) --go-grpc_opt=paths=source_relative \
		inventory/v1/inventory.proto

gen-user:
	mkdir -p $(OUT)
	$(PROTOC) \
		--go_out=$(OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT) --go-grpc_opt=paths=source_relative \
		user/v1/user.proto

clean:
	rm -rf $(OUT)
