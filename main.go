package main

import (
	"log"
	"net/http"
	"phikhanh/config"
	"phikhanh/middlewares"
	"phikhanh/routes"
	"phikhanh/utils"
	"time"

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
// @description Provide your JWT token for the Authorization header (with or without the "Bearer " prefix).
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

	// Suppress Chrome DevTools auto-discovery request
	router.GET("/.well-known/appspecific/com.chrome.devtools.json", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	// Load HTML templates
	router.SetHTMLTemplate(utils.LoadTemplates("templates"))

	// Serve static files
	router.Static("/assets", "./assets")

	// Thiết lập routes
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/admin/login")
	})

	routes.SetupUserRoutes(router)
	routes.SetupAdminRoutes(router)

	// Khởi động server
	serverAddr := ":" + config.AppConfig.ServerPort
	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	log.Printf("Admin: http://localhost:%s/admin/login", config.AppConfig.ServerPort)
	log.Printf("Swagger UI: http://localhost:%s/docs/index.html", config.AppConfig.ServerPort)

	// Dọn dẹp upload records định kỳ mỗi giờ
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			middlewares.CleanupUploadRecords()
		}
	}()

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
