package grpc

import (
	"authentication/internal/token"
	proto "authentication/internal/token/proto/token"
	"context"
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
	response := th.tokenService.IsTokenValid(c, req.Token)

	return &proto.CheckTokenResponse{
		Valid: response,
	}, nil
}
