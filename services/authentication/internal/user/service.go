package user

import (
	"authentication/internal/platform/encryption"
	"context"
)

type UserRepository interface {
	Create(c context.Context, user *User) (*User, error)
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
