package handler

import (
	"context"
	"go-multitenant/internal/model"
	"go-multitenant/internal/service"
	pb "go-multitenant/proto/userpb"
)

type GRPCServer struct {
	pb.UnimplementedUserServiceServer
	UserService *service.UserService
}

func NewGRPCServer(us *service.UserService) *GRPCServer { return &GRPCServer{UserService: us} }
func (s *GRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	out, err := s.UserService.GetUser(ctx, req.ClientId, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	u := out.(*model.User)
	return &pb.GetUserResponse{UserId: int64(u.ID), Name: u.Name}, nil
}
