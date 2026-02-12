package config

import (
	"log"
	"phikhanh/models"
)

// Chạy auto migration cho các models
func RunMigrations() {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.Service{},
		&models.Application{},
		&models.Attachment{},
		&models.ApplicationHistory{},
		&models.Notification{},
		&models.SystemLog{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}
