package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	inventoryv1 "github.com/adopabianko/commerce/proto/gen/inventory/v1"
)

func New(addr string) (inventoryv1.InventoryServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { return nil, err }
	return inventoryv1.NewInventoryServiceClient(conn), nil
}
