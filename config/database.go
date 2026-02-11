package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Kết nối đến PostgreSQL database
func ConnectDatabase() {
	var err error
	dsn := AppConfig.GetDSN()

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

// Trả về instance database hiện tại
func GetDB() *gorm.DB {
	return DB
}
