package auth

import (
	"context"
	tokenPb "core/internal/proto/token"
	userPb "core/internal/proto/user"
)

type grpcAdapter struct {
	userClient  userPb.UserServiceClient
	tokenClient tokenPb.TokenServiceClient
}

func NewGRPCAdapter(uc userPb.UserServiceClient, tc tokenPb.TokenServiceClient) *grpcAdapter {
	return &grpcAdapter{
		userClient:  uc,
		tokenClient: tc,
	}
}

func (a *grpcAdapter) Register(c context.Context, phoneNumber, name, password string) (*AuthResponse, error) {
	req := userPb.RegisterRequest{
		Name:        name,
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	res, err := a.userClient.RegisterUser(c, &req)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		Name:         res.Name,
		UUID:         res.Uuid,
	}, nil
}

func (a *grpcAdapter) Login(c context.Context, phoneNumber, password string) (*AuthResponse, error) {
	req := userPb.LoginRequest{
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	res, err := a.userClient.Login(c, &req)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		Name:         res.Name,
		UUID:         res.Uuid,
	}, err
}

func (a *grpcAdapter) Logout(c context.Context, token string) error {
	req := userPb.LogoutRequset{
		Token: token,
	}

	_, err := a.userClient.Logout(c, &req)
	if err != nil {
		return err
	}

	return nil
}

func (a *grpcAdapter) IsTokenValid(c context.Context, token string) (*CheckTokenResponse, error) {
	req := tokenPb.CheckTokenRequest{
		Token: token,
	}

	res, err := a.tokenClient.IsTokenValid(c, &req)
	if err != nil {
		return nil, err
	}

	return &CheckTokenResponse{
		Valid: res.Valid,
		UUID:  res.Uuid,
	}, nil
}

func (a *grpcAdapter) Refresh(c context.Context, refreshToken string) (*RefreshResponse, error) {
	req := tokenPb.RefreshRequest{
		RefreshToken: refreshToken,
	}

	res, err := a.tokenClient.Refresh(c, &req)
	if err != nil {
		return nil, err
	}

	return &RefreshResponse{
		RefreshToken: res.RefreshToken,
		AccessToken:  res.AcessToken,
	}, nil
}
