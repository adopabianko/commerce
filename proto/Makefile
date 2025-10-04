PROTOC ?= protoc
OUT := ./gen

all: clean gen

gen:
	mkdir -p $(OUT)
	$(PROTOC) --go_out=$(OUT) --go_opt=paths=source_relative 		--go-grpc_out=$(OUT) --go-grpc_opt=paths=source_relative 		inventory/v1/inventory.proto

clean:
	rm -rf $(OUT)
