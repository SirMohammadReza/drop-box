package grpc

import (
	"authentication/internal/user"
	proto "authentication/internal/user/proto/user"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer
	userProvider user.Provider
}

func NewUserHandler(up user.Provider) *UserHandler {
	return &UserHandler{
		userProvider: up,
	}
}

func (uh *UserHandler) Login(c context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	response, err := uh.userProvider.Login(c, req.PhoneNumber, req.Password)
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
	if err := uh.userProvider.Logout(c, req.Token); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (uh *UserHandler) RegisterUser(c context.Context, req *proto.RegisterRequest) (*proto.AuthResponse, error) {
	response, err := uh.userProvider.RegisterUser(c, req.Name, req.PhoneNumber, req.Password)
	if err != nil {
		return nil, err
	}

	return &proto.AuthResponse{
		RefreshToken: response.RefreshToken,
		AccessToken:  response.AccessToken,
		Name:         response.Name,
	}, nil
}
