package main

import (
	"log"
	"phikhanh/config"
	"phikhanh/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Khởi tạo cấu hình
	config.LoadConfig()

	// Kết nối database
	config.ConnectDatabase()

	// Chạy migration
	config.RunMigrations()

	// Khởi tạo Gin router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/**/*")

	// Serve static files
	router.Static("/assets", "./assets")

	// Thiết lập routes
	routes.SetupRoutes(router)

	// Khởi động server
	serverAddr := ":" + config.AppConfig.ServerPort
	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
