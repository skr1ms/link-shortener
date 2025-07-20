package migrations

import (
	"fmt"
	"linkshortener/configs"
	"linkshortener/internal/link"
	"linkshortener/internal/stats"
	"linkshortener/internal/user"
	"linkshortener/pkg/db"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func waitForDB(config *configs.Config) {
	for {
		db, err := gorm.Open(postgres.Open(config.DB.DSN), &gorm.Config{})
		if err == nil {
			if sqlDB, err := db.DB(); err == nil {
				if err := sqlDB.Ping(); err == nil {
					fmt.Println("Database connected!")
					sqlDB.Close()
					break
				}
				sqlDB.Close()
			}
		}
		fmt.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}
}

func RunMigrations(config *configs.Config) *db.Db {
	waitForDB(config)

	database := db.NewDb(config)

	err := database.AutoMigrate(&link.Link{}, &user.User{}, &stats.Stats{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database migrations completed successfully!")
	return database
}
