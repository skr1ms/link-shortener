package db

import (
	"linkshortener/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(config *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(config.DB.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return &Db{db}
}