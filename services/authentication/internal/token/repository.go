package token

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (pr *PostgresRepository) StoreRefreshToken(c context.Context, token *Token) error {
	res := pr.db.WithContext(c).Create(&token)
	if res.Error != nil {
		return fmt.Errorf("Cound not store token: %w", res.Error)
	}

	return nil
}

func (pr *PostgresRepository) GetToken(c context.Context, userID uint, token string) (*Token, error) {
	var t Token

	res := pr.db.WithContext(c).Where("user_id = ?", userID).Where("token = ?", token).First(&t)
	if res.Error != nil {
		return nil, fmt.Errorf("Token not found: %w", res.Error)
	}

	return &t, nil
}

func (pr *PostgresRepository) RevokeToken(c context.Context, token string) error {
	return pr.db.WithContext(c).Unscoped().Where("token = ?", token).Delete(&Token{}).Error
}
