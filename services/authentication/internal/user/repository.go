package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(c context.Context, user *User) (*User, error) {
	res := r.db.WithContext(c).Create(user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *Repository) FindByPhoneNumber(c context.Context, phoneNumber string) (*User, error) {
	var user User
	res := r.db.Where("phone_number = ?", phoneNumber).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}
