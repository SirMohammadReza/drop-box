package grpc

import (
	"authentication/internal/user"
	proto "authentication/internal/user/proto/user"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer
	userService *user.UserService
}

func NewUserHandler(us *user.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (uh *UserHandler) Login(c context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	response, err := uh.userService.Login(c, req.PhoneNumber, req.Password)
	if err != nil {
		return nil, err
	}

	return &proto.AuthResponse{
		RefreshToken: response.RefreshToken,
		AccessToken:  response.AccessToken,
		Name:         response.Name,
	}, nil
}

func (uh *UserHandler) Logout(c context.Context, req *proto.LogoutRequset) (*emptypb.Empty, error) {
	if err := uh.userService.Logout(c, req.Token); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (uh *UserHandler) RegisterUser(c context.Context, req *proto.RegisterRequest) (*proto.AuthResponse, error) {
	response, err := uh.userService.RegisterUser(c, req.Name, req.PhoneNumber, req.Password)
	if err != nil {
		return nil, err
	}

	return &proto.AuthResponse{
		RefreshToken: response.RefreshToken,
		AccessToken:  response.AccessToken,
		Name:         response.Name,
	}, nil
}
