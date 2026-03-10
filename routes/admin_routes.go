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
	db := config.GetDB()

	// Activity Logs (shared across controllers)
	activityLogRepo := adminRepo.NewActivityLogRepository(db)
	activityLogService := adminSvc.NewActivityLogService(activityLogRepo)
	activityLogCtrl := admin.NewActivityLogController(activityLogService)

	// Auth
	authRepo := adminRepo.NewAdminRepository(db)
	authService := adminSvc.NewAuthService(authRepo)
	authController := admin.NewAuthController(authService, activityLogService)

	// Services CRUD
	svcRepo := adminRepo.NewServiceRepository(db)
	svcService := adminSvc.NewServiceAdminService(svcRepo)
	serviceController := admin.NewServiceController(svcService, activityLogService)

	// Departments CRUD
	deptRepo := adminRepo.NewDepartmentRepository(db)
	deptService := adminSvc.NewDepartmentService(deptRepo)
	departmentController := admin.NewDepartmentController(deptService, activityLogService)

	// Applications
	appRepo := adminRepo.NewApplicationRepository(db)
	appService := adminSvc.NewApplicationAdminService(appRepo)
	applicationController := admin.NewApplicationController(appService, activityLogService)

	// Users
	userRepo := adminRepo.NewUserRepository(db)
	userService := adminSvc.NewUserService(userRepo)
	userController := admin.NewUserController(userService, activityLogService)

	// Dashboard
	dashboardController := admin.NewDashboardController(appService, svcService, deptService)

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/login", authController.ShowLogin)
		adminGroup.POST("/login", authController.ProcessLogin)
		adminGroup.POST("/logout", authController.ProcessLogout)

		protected := adminGroup.Group("")
		protected.Use(middlewares.AdminAuthMiddleware())
		{
			protected.GET("/dashboard", dashboardController.ShowDashboard)

			applications := protected.Group("/applications")
			applications.Use(middlewares.RequireRole("admin", "manager", "staff"))
			{
				applications.GET("", applicationController.List)
				applications.GET("/:id", applicationController.ShowDetail)
				applications.POST("/:id/process", applicationController.Process)
			}

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

			activityLogs := protected.Group("/activity-logs")
			activityLogs.Use(middlewares.RequireRole("admin"))
			{
				activityLogs.GET("", activityLogCtrl.List)
				activityLogs.POST("/cleanup", activityLogCtrl.Cleanup)
			}
		}
	}
}
