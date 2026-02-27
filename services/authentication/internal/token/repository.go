package token

import (
	"context"
	"fmt"

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

func (r *Repository) StoreRefreshToken(c context.Context, token *Token) error {
	res := r.db.WithContext(c).Create(&token)
	if res.Error != nil {
		return fmt.Errorf("Cound not store token: %w", res.Error)
	}

	return nil
}

func (r *Repository) GetToken(c context.Context, userID uint, token string) (*Token, error) {
	var t Token

	res := r.db.WithContext(c).Where("user_id = ?", userID).Where("token = ?", token).First(&t)
	if res.Error != nil {
		return nil, fmt.Errorf("Token not found: %w", res.Error)
	}

	return &t, nil
}

func (r *Repository) RevokeToken(c context.Context, token string) error {
	var t Token
	res := r.db.WithContext(c).Unscoped().Where("token = ?", token).First(&t).Delete(&t)
	if res.Error != nil {
		return fmt.Errorf("could nod delete refresh token: %w", res.Error)
	}

	return nil
}
