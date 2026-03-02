package routes

import (
	"phikhanh/config"
	"phikhanh/controllers/admin"
	"phikhanh/middlewares"
	adminRepo "phikhanh/repositories/admin"
	adminSvc "phikhanh/services/admin"

	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes - Thiết lập routes cho Admin
func SetupAdminRoutes(router *gin.Engine) {
	repo := adminRepo.NewAdminRepository(config.GetDB())
	authService := adminSvc.NewAuthService(repo)

	authController := admin.NewAuthController(authService)
	dashboardController := admin.NewDashboardController()
	userController := admin.NewUserController()
	serviceController := admin.NewServiceController()
	departmentController := admin.NewDepartmentController()
	applicationController := admin.NewApplicationController()
	activityLogController := admin.NewActivityLogController()

	adminGroup := router.Group("/admin")
	{
		// Auth routes (public)
		adminGroup.GET("/login", authController.ShowLogin)
		adminGroup.POST("/login", authController.ProcessLogin)
		adminGroup.POST("/logout", authController.ProcessLogout)

		// Protected routes
		protected := adminGroup.Group("")
		protected.Use(middlewares.AdminAuthMiddleware())
		{
			protected.GET("/dashboard", dashboardController.ShowDashboard)
			protected.GET("/users", userController.ShowUsers)
			protected.GET("/services", serviceController.ShowServices)
			protected.GET("/departments", departmentController.ShowDepartments)
			protected.GET("/applications", applicationController.ShowApplications)
			protected.GET("/activity-logs", activityLogController.ShowActivityLogs)
		}
	}
}
