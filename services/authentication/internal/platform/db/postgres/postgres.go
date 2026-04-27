package postgres

import (
	"authentication/internal/config"
	"fmt"
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
		dsn := getDsn()
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Could not create database: %s", err.Error())
		}
		_db = db
	})

	return _db
}

func getDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.Config("PG_HOST"), config.Config("PG_USER"), config.Config("PG_PASSWORD"), config.Config("PG_DBNAME"), config.Config("PG_PORT"), config.Config("PG_SSLMODE"))
}
