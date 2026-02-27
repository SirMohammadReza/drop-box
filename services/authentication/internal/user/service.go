package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Provider interface {
	RegisterUser(c context.Context, name, phoneNumber, password string) (*AuthResponse, error)
	GetUserInfo(c context.Context, phoneNumber string) (*User, error)
	Login(c context.Context, phoneNumber, password string) (*AuthResponse, error)
	Logout(c context.Context, token string) error
}

type UserRepository interface {
	Create(c context.Context, user *User) (*User, error)
	FindByPhoneNumber(c context.Context, phoneNumber string) (*User, error)
}

type TokenProvider interface {
	GenerateTokenPair(c context.Context, userID uint, uuid uuid.UUID) (string, string, error)
	DeleteToken(c context.Context, token string) error
}

type PasswordManager interface {
	HashPassword(password string) (string, error)
	CheckHash(password, hashValue string) bool
}

type UserService struct {
	userRepository  UserRepository
	passwordManager PasswordManager
	tokenProvider   TokenProvider
}

type AuthResponse struct {
	AccessToken  string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	Name         string    `json:"name"`
	Uuid         uuid.UUID `json:"uuid"`
}

func NewUserService(ur UserRepository, pm PasswordManager, tp TokenProvider) Provider {
	return &UserService{
		userRepository:  ur,
		passwordManager: pm,
		tokenProvider:   tp,
	}
}

func (us *UserService) RegisterUser(c context.Context, name, phoneNumber, password string) (*AuthResponse, error) {
	hashPassword, err := us.passwordManager.HashPassword(password)
	if err != nil {
		return nil, err
	}
	newUuid := uuid.New()
	newUser := User{
		Name:        name,
		Uuid:        newUuid,
		PhoneNumber: phoneNumber,
		Password:    hashPassword,
	}

	user, err := us.userRepository.Create(c, &newUser)

	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := us.tokenProvider.GenerateTokenPair(c, user.ID, newUuid)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         user.Name,
		Uuid:         newUuid,
	}, nil
}

func (us *UserService) GetUserInfo(c context.Context, phoneNumber string) (*User, error) {
	user, err := us.userRepository.FindByPhoneNumber(c, phoneNumber)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Login(c context.Context, phoneNumber, password string) (*AuthResponse, error) {
	user, err := us.userRepository.FindByPhoneNumber(c, phoneNumber)
	if err != nil {
		return nil, err
	}

	valid := us.passwordManager.CheckHash(password, user.Password)
	if !valid {
		return nil, errors.New("Incorrect password")
	}

	accessToken, refreshToken, err := us.tokenProvider.GenerateTokenPair(c, user.ID, user.Uuid)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         user.Name,
		Uuid:         user.Uuid,
	}, nil
}

func (us *UserService) Logout(c context.Context, token string) error {
	return us.tokenProvider.DeleteToken(c, token)
}
