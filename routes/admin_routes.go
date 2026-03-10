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

	// Services CRUD
	svcRepo := adminRepo.NewServiceRepository(config.GetDB())
	svcService := adminSvc.NewServiceAdminService(svcRepo)
	serviceController := admin.NewServiceController(svcService)

	// Departments CRUD
	deptRepo := adminRepo.NewDepartmentRepository(config.GetDB())
	deptService := adminSvc.NewDepartmentService(deptRepo)
	departmentController := admin.NewDepartmentController(deptService)

	// Applications
	appRepo := adminRepo.NewApplicationRepository(config.GetDB())
	appService := adminSvc.NewApplicationAdminService(appRepo)
	applicationController := admin.NewApplicationController(appService)

	// Users
	userRepo := adminRepo.NewUserRepository(config.GetDB())
	userService := adminSvc.NewUserService(userRepo)
	userController := admin.NewUserController(userService)

	// Activity Logs
	activityLogController := admin.NewActivityLogController()

	// Dashboard
	dashboardController := admin.NewDashboardController(appService, svcService, deptService)

	// Setup routes
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/login", authController.ShowLogin)
		adminGroup.POST("/login", authController.ProcessLogin)
		adminGroup.POST("/logout", authController.ProcessLogout)

		protected := adminGroup.Group("")
		protected.Use(middlewares.AdminAuthMiddleware())
		{
			protected.GET("/dashboard", dashboardController.ShowDashboard)

			// Applications - All roles (admin, manager, staff)
			applications := protected.Group("/applications")
			applications.Use(middlewares.RequireRole("admin", "manager", "staff"))
			{
				applications.GET("", applicationController.List)
				applications.GET("/:id", applicationController.ShowDetail)
				applications.POST("/:id/process", applicationController.Process)
			}

			// Services - Admin & Manager only
			services := protected.Group("/services")
			services.Use(middlewares.RequireRole("admin", "manager"))
			{
				services.GET("", serviceController.List)
				services.GET("/create", serviceController.CreateForm)
				services.POST("/create", serviceController.CreateSave)
				services.GET("/:id/edit", serviceController.EditForm)
				services.POST("/:id/edit", serviceController.EditSave)
				services.POST("/:id/delete", serviceController.Delete)
			}

			// Departments - Admin only
			departments := protected.Group("/departments")
			departments.Use(middlewares.RequireRole("admin"))
			{
				departments.GET("", departmentController.List)
				departments.GET("/create", departmentController.CreateForm)
				departments.POST("/create", departmentController.CreateSave)
				departments.GET("/:id", departmentController.Detail)
				departments.GET("/:id/edit", departmentController.EditForm)
				departments.POST("/:id/edit", departmentController.EditSave)
				departments.POST("/:id/delete", departmentController.Delete)
			}

			// Users - Admin only
			users := protected.Group("/users")
			users.Use(middlewares.RequireRole("admin"))
			{
				users.GET("", userController.List)
				users.GET("/create", userController.CreateForm)
				users.POST("/create", userController.CreateSave)
				users.GET("/:id", userController.ShowDetail)
				users.GET("/:id/edit", userController.EditForm)
				users.POST("/:id/edit", userController.EditSave)
				users.POST("/:id/delete", userController.Delete)
			}

			// Activity Logs - Admin only
			activityLogs := protected.Group("/activity-logs")
			activityLogs.Use(middlewares.RequireRole("admin"))
			{
				activityLogs.GET("", activityLogController.ShowActivityLogs)
			}
		}
	}
}
