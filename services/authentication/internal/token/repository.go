package token

import (
	"context"

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
		Toekn:  token,
	}

	return pr.db.WithContext(c).Create(t).Error
}

func (pr *PostgresRepository) RevokeToken(c context.Context, token string) error {
	return pr.db.WithContext(c).Unscoped().Where("token = ?", token).Delete(&Token{}).Error
}

func (pr *PostgresRepository) RevokeTokenByUserID(c context.Context, userID uint) error {
	return pr.db.WithContext(c).Unscoped().Where("user_id = ?", userID).Delete(&Token{}).Error
}
