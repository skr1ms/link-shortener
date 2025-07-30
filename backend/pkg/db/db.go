package db

import (
	"linkshortener/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(config *config.Config) *Db {
	db, err := gorm.Open(postgres.Open(config.DB.URL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return &Db{db}
}