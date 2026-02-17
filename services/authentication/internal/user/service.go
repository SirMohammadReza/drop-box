package user

import (
	"authentication/internal/platform/encryption"
	"context"
)

type UserRepository interface {
	Create(c context.Context, user User) error
	FindByPhoneNumber(c context.Context, phoneNumber string) (*User, error)
}

type UserService struct {
	userRepo        UserRepository
	passwordManager encryption.EncryptionService
}

func NewUserService(ur UserRepository, pm encryption.EncryptionService) *UserService {
	return &UserService{
		userRepo:        ur,
		passwordManager: pm,
	}
}

func (us *UserService) StoreNewUser(c context.Context, name, phoneNumber, password string) error {
	hashPassword, err := us.passwordManager.HashPassword(password)
	if err != nil {
		return err
	}

	newUser := User{
		Name:        name,
		PhoneNumber: phoneNumber,
		Password:    hashPassword,
	}

	return us.userRepo.Create(c, newUser)
}

func (us *UserService) GetUserInfo(c context.Context, phoneNumber string) (*User, error) {
	user, err := us.userRepo.FindByPhoneNumber(c, phoneNumber)
	if err != nil {
		return nil, err
	}

	return user, nil
}
