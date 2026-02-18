package postgres

import (
	"authentication/internal/user"
	"context"

	"gorm.io/gorm"
)

type PostgresDB struct {
	DB *gorm.DB
}

func (p *PostgresDB) Create(c context.Context, user *user.User) (*user.User, error) {
	result := p.DB.WithContext(c).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (p *PostgresDB) FindByPhoneNumber(c context.Context, phoneNumber string) (*user.User, error) {
	var user user.User

	result := p.DB.WithContext(c).Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
