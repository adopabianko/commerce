package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	userv1 "github.com/adopabianko/commerce/proto/gen/user/v1"
	"google.golang.org/grpc"
)

type GRPCAuthClient struct {
	client userv1.UserServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCAuthClient() (*GRPCAuthClient, error) {
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = "user-service:50052" // default docker compose address
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed connect user-service grpc: %w", err)
	}

	client := userv1.NewUserServiceClient(conn)
	return &GRPCAuthClient{client: client, conn: conn}, nil
}

func (g *GRPCAuthClient) ValidateToken(token string) (uint32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := g.client.Validate(ctx, &userv1.ValidateRequest{Token: token})
	if err != nil {
		return 0, err
	}
	if !res.Valid {
		return 0, fmt.Errorf("invalid token: %s", res.Error)
	}
	return res.UserId, nil
}

func (g *GRPCAuthClient) Close() error {
	if g.conn == nil {
		return nil
	}
	return g.conn.Close()
}
