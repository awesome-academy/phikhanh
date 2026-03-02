package routes

import (
	"phikhanh/config"
	"phikhanh/controllers/admin"
	"phikhanh/middlewares"
	adminRepo "phikhanh/repositories/admin"
	adminSvc "phikhanh/services/admin"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.Engine) {
	// Auth
	authRepo := adminRepo.NewAdminRepository(config.GetDB())
	authService := adminSvc.NewAuthService(authRepo)
	authController := admin.NewAuthController(authService)

	// Service CRUD
	svcRepo := adminRepo.NewServiceRepository(config.GetDB())
	svcService := adminSvc.NewServiceAdminService(svcRepo)
	serviceController := admin.NewServiceController(svcService)

	// Other controllers
	dashboardController := admin.NewDashboardController()
	userController := admin.NewUserController()
	departmentController := admin.NewDepartmentController()
	applicationController := admin.NewApplicationController()
	activityLogController := admin.NewActivityLogController()

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/login", authController.ShowLogin)
		adminGroup.POST("/login", authController.ProcessLogin)
		adminGroup.POST("/logout", authController.ProcessLogout)

		protected := adminGroup.Group("")
		protected.Use(middlewares.AdminAuthMiddleware())
		{
			protected.GET("/dashboard", dashboardController.ShowDashboard)
			protected.GET("/users", userController.ShowUsers)
			protected.GET("/departments", departmentController.ShowDepartments)
			protected.GET("/applications", applicationController.ShowApplications)
			protected.GET("/activity-logs", activityLogController.ShowActivityLogs)

			// Services CRUD
			services := protected.Group("/services")
			{
				services.GET("", serviceController.List)
				services.GET("/create", serviceController.CreateForm)
				services.POST("/create", serviceController.CreateSave)
				services.GET("/:id", serviceController.Detail)
				services.GET("/:id/edit", serviceController.EditForm)
				services.POST("/:id/edit", serviceController.EditSave)
				services.POST("/:id/delete", serviceController.Delete)
			}
		}
	}
}
