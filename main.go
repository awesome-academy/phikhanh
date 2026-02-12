package main

import (
	"log"
	"phikhanh/config"
	"phikhanh/routes"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"

	_ "phikhanh/docs"
)

// @title           Public Service Management API
// @version         1.0
// @description     API cho hệ thống quản lý dịch vụ công
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Khởi tạo cấu hình
	config.LoadConfig()

	// Kết nối database
	config.ConnectDatabase()

	// Chạy migration
	config.RunMigrations()

	// Đăng ký custom validators
	utils.RegisterCustomValidators()

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
	log.Printf("Swagger UI: http://localhost:%s/docs/index.html", config.AppConfig.ServerPort)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
