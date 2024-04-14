package grpc

import (
	"context"
	"net"

	"github.com/Alieksieiev0/feed-service/api/proto"
	"github.com/Alieksieiev0/feed-service/internal/models"
	"github.com/Alieksieiev0/feed-service/internal/services"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	addr string
}

func NewServer(addr string) *GRPCServer {
	return &GRPCServer{
		addr: addr,
	}
}

func (s *GRPCServer) Start(service services.UserService) error {
	grpcUserService := NewGRPCUserServiceServer(service)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(server, grpcUserService)

	return server.Serve(ln)
}

type GRPCUserServiceServer struct {
	service services.UserService
	proto.UnimplementedUserServiceServer
}

func NewGRPCUserServiceServer(service services.UserService) proto.UserServiceServer {
	return &GRPCUserServiceServer{
		service: service,
	}
}

func (us *GRPCUserServiceServer) GetByUsername(
	c context.Context,
	req *proto.UsernameRequest,
) (*proto.UserResponse, error) {
	user, err := us.service.GetByUsername(c, req.Username)
	if err != nil {
		return nil, err
	}

	resp := &proto.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
	return resp, err
}

func (us *GRPCUserServiceServer) Save(
	c context.Context,
	req *proto.UserRequest,
) (*proto.SaveResponse, error) {
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	err := us.service.Save(c, user)
	if err != nil {
		return nil, err
	}

	return &proto.SaveResponse{}, nil
}
