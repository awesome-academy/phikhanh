package config

import (
	"log"
)

// Chạy auto migration cho các models
func RunMigrations() {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate()

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}
