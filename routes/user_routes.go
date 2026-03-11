package routes

import (
	"phikhanh/config"
	userCtrl "phikhanh/controllers/user"
	"phikhanh/middlewares"
	"phikhanh/repositories"
	userRepo "phikhanh/repositories/user"
	"phikhanh/services"
	userSvc "phikhanh/services/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const maxUploadBodySize = int64(11 << 20) // 11MB

// Thiết lập routes cho User API (JSON)
func SetupUserRoutes(router *gin.Engine) {
	db := config.GetDB()

	// Health check
	healthController := userCtrl.NewHealthController()

	// Auth
	authRepo := userRepo.NewAuthRepository(db)
	authService := userSvc.NewAuthService(authRepo)
	authController := userCtrl.NewAuthController(authService)

	// Profile
	profileRepo := userRepo.NewProfileRepository(db)
	profileService := userSvc.NewProfileService(profileRepo)
	profileController := userCtrl.NewProfileController(profileService)

	// Service
	serviceRepo := userRepo.NewServiceRepository(db)
	serviceService := userSvc.NewServiceService(serviceRepo)
	serviceController := userCtrl.NewServiceController(serviceService)

	// Upload
	uploadController := userCtrl.NewUploadController()

	// Application
	appRepo := userRepo.NewApplicationRepository(db)
	appService := userSvc.NewApplicationService(appRepo)
	appController := userCtrl.NewApplicationController(appService)

	// Notifications
	notifRepo := repositories.NewNotificationRepository(db)
	notifService := services.NewNotificationService(notifRepo)
	notifController := userCtrl.NewNotificationController(notifService)

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

		// Service routes (public)
		services := api.Group("/services")
		{
			services.GET("", serviceController.GetServiceList)
			services.GET("/:id", serviceController.GetServiceDetail)
		}

		// Upload routes (protected + rate limited)
		api.POST("/upload",
			middlewares.AuthMiddleware(),
			middlewares.UploadRateLimitMiddleware(),
			uploadController.UploadFile,
		)

		// Application routes (protected)
		applications := api.Group("/applications")
		applications.Use(middlewares.AuthMiddleware())
		{
			applications.POST("", appController.SubmitApplication)
			applications.GET("/me", appController.GetMyApplications)
			applications.POST("/:id/supplement", appController.SupplementApplication)
		}

		// Notifications routes (protected)
		notifications := api.Group("/notifications")
		notifications.Use(middlewares.AuthMiddleware())
		{
			notifications.GET("", notifController.GetNotifications)
			notifications.PATCH("/read-all", notifController.MarkAllAsRead)
			notifications.PATCH("/:id/read", notifController.MarkAsRead)
		}
	}
}
