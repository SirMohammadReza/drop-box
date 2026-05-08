package auth

import "context"

type Provider interface {
	Register(c context.Context, phoneNumber, name, password string) (*AuthResponse, error)
	Login(c context.Context, phoneNumber, password string) (*AuthResponse, error)
	Logout(c context.Context, token string) error

	IsTokenValid(c context.Context, token string) (*CheckTokenResponse, error)
	Refresh(c context.Context, refreshToken string) (*RefreshResponse, error)
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
}

type CheckTokenResponse struct {
	Valid bool   `json:"valid"`
	UUID  string `json:"uuid"`
}

type RefreshResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
