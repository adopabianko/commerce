package grpcsvr

import (
	"context"

	userv1 "github.com/adopabianko/commerce/proto/gen/user/v1"
	"github.com/adopabianko/commerce/user-service/internal/usecase"
)

type Server struct {
	userv1.UnimplementedUserServiceServer
	svc *usecase.Service
}

func NewServer(s *usecase.Service) *Server { return &Server{svc: s} }

func (s *Server) Validate(ctx context.Context, req *userv1.ValidateRequest) (*userv1.ValidateResponse, error) {
	uid, err := s.svc.ValidateToken(req.Token)
	if err != nil {
		return &userv1.ValidateResponse{Valid: false, Error: err.Error()}, nil
	}
	return &userv1.ValidateResponse{Valid: true, UserId: uint32(uid)}, nil
}
