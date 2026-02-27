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
	tokenProvider token.Provider
}

func NewTokenHandler(tp token.Provider) *TokenHandler {
	return &TokenHandler{
		tokenProvider: tp,
	}
}

func (th *TokenHandler) IsTokenValid(c context.Context, req *proto.CheckTokenRequest) (*proto.CheckTokenResponse, error) {
	_, err := th.tokenProvider.ValidateToken(c, req.Token, token.TokenAccessType)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "could not validate access token")
	}

	return &proto.CheckTokenResponse{
		Valid: true,
	}, nil
}

func (th *TokenHandler) Refresh(c context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	accessToken, refreshToken, err := th.tokenProvider.RefreshTokens(c, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "could not refresh token")
	}

	return &proto.RefreshResponse{
		AcessToken:   accessToken,
		RefreshToken: refreshToken,
	}, nil
}
