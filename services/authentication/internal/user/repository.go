package user

import (
	"context"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (pr *PostgresRepository) Create(c context.Context, user *User) (*User, error) {
	res := pr.DB.WithContext(c).Create(user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (pr *PostgresRepository) FindByPhoneNumber(c context.Context, phoneNumber string) (*User, error) {
	var user User
	res := pr.DB.Where("phone_number = ?", phoneNumber).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}
