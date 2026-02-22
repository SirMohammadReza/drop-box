package token

import (
	"context"
	"errors"

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

func (pr *PostgresRepository) StoreRefreshToken(c context.Context, userID uint, token string) error {
	t := Token{
		UserID: userID,
		Token:  token,
	}

	return pr.db.WithContext(c).Create(&t).Error
}

func (pr *PostgresRepository) IsTokenExists(c context.Context, token string) bool {
	err := pr.db.WithContext(c).First(&Token{}, "token = ?", token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

func (pr *PostgresRepository) RevokeToken(c context.Context, token string) error {
	return pr.db.WithContext(c).Unscoped().Where("token = ?", token).Delete(&Token{}).Error
}

func (pr *PostgresRepository) RevokeTokenByUserID(c context.Context, userID uint) error {
	return pr.db.WithContext(c).Unscoped().Where("user_id = ?", userID).Delete(&Token{}).Error
}
