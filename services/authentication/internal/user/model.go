package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `gorm:"type:varchar(30),not null" json:"username"`
	PhoneNumber string `gorm:"type:varchar(13),not null,unique" json:"phone_number"`
	Password    string `gorm:"type:varchar(255),not null" json:"_"`
}
