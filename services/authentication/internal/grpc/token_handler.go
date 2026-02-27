package grpc

import (
	"authentication/internal/token"
	proto "authentication/internal/token/proto/token"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenHandler struct {
	proto.UnimplementedTokenServiceServer
	tokenService *token.TokenService
}

func NewTokenHandler(ts *token.TokenService) *TokenHandler {
	return &TokenHandler{
		tokenService: ts,
	}
}

func (th *TokenHandler) IsTokenValid(c context.Context, req *proto.CheckTokenRequest) (*proto.CheckTokenResponse, error) {
	_, err := th.tokenService.ValidateToken(c, req.Token, token.TokenAccessType)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "could not validate access token")
	}

	return &proto.CheckTokenResponse{
		Valid: true,
	}, nil
}

func (th *TokenHandler) Refresh(c context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	accessToken, refreshToken, err := th.tokenService.RefreshTokens(c, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "could not refresh token")
	}

	return &proto.RefreshResponse{
		AcessToken:   accessToken,
		RefreshToken: refreshToken,
	}, nil
}
