package user

import (
	"authentication/internal/platform/encryption"
	"authentication/internal/token"
	"context"
	"errors"
)

type UserRepository interface {
	Create(c context.Context, user *User) (*User, error)
	FindByPhoneNumber(c context.Context, phoneNumber string) (*User, error)
}

type UserService struct {
	userRepo        UserRepository
	passwordManager encryption.EncryptionService
	tokenService    *token.TokenService
}

type AuthResponse struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func NewUserService(ur UserRepository, pm encryption.EncryptionService, ts token.TokenService) *UserService {
	return &UserService{
		userRepo:        ur,
		passwordManager: pm,
		tokenService:    &ts,
	}
}

func (us *UserService) StoreNewUser(c context.Context, name, phoneNumber, password string) (*User, error) {
	hashPassword, err := us.passwordManager.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := User{
		Name:        name,
		PhoneNumber: phoneNumber,
		Password:    hashPassword,
	}

	user, err := us.userRepo.Create(c, &newUser)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserInfo(c context.Context, phoneNumber string) (*User, error) {
	user, err := us.userRepo.FindByPhoneNumber(c, phoneNumber)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Login(c context.Context, phoneNumber, password string) (*AuthResponse, error) {
	user, err := us.userRepo.FindByPhoneNumber(c, phoneNumber)
	if err != nil {
		return nil, err
	}

	valid := us.passwordManager.CheckHash(password, user.Password)
	if !valid {
		return nil, errors.New("Incorrect password")
	}

	accessToken, refreshToken, err := us.tokenService.GenerateTokenPair(c, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (us *UserService) Logout(c context.Context, token string) error {
	return us.tokenService.DeleteToken(c, token)
}
