package postgres

import (
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	_db  *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		dsn := "host=postgres user=admin password=4321 dbname=auth port=5432 sslmode=disable TimeZone=Asia/Tehran"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Could not create database: %s", err.Error())
		}
		_db = db
	})

	return _db
}
