package token

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Token  string `gorm:"varchar(255);not null" json:"token"`
}
