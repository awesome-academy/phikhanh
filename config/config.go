package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Cấu trúc lưu trữ các biến môi trường
type Config struct {
	ServerPort string
	GinMode    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

var AppConfig *Config

// Đọc file .env và khởi tạo cấu hình
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "public_service_management"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	log.Println("Configuration loaded successfully")
}

// Lấy giá trị biến môi trường hoặc trả về giá trị mặc định
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Tạo chuỗi kết nối PostgreSQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}
