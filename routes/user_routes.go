package routes

import (
	"phikhanh/config"
	userCtrl "phikhanh/controllers/user"
	"phikhanh/middlewares"
	userRepo "phikhanh/repositories/user"
	userSvc "phikhanh/services/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Thiết lập routes cho User API (JSON)
func SetupUserRoutes(router *gin.Engine) {
	// Health check
	healthController := userCtrl.NewHealthController()

	// Auth
	authRepo := userRepo.NewAuthRepository(config.GetDB())
	authService := userSvc.NewAuthService(authRepo)
	authController := userCtrl.NewAuthController(authService)

	// Profile
	profileRepo := userRepo.NewProfileRepository(config.GetDB())
	profileService := userSvc.NewProfileService(profileRepo)
	profileController := userCtrl.NewProfileController(profileService)

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Group API v1
	api := router.Group("/api/v1")
	{
		api.GET("/health", healthController.CheckHealth)

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/logout", middlewares.AuthMiddleware(), authController.Logout)
		}

		// Profile routes (protected)
		profile := api.Group("/profile")
		profile.Use(middlewares.AuthMiddleware())
		{
			profile.GET("", profileController.GetProfile)
			profile.PUT("", profileController.UpdateProfile)
		}
	}
}
