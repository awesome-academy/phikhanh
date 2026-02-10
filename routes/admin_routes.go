package routes

import (
	adminCtrl "phikhanh/controllers/admin"

	"github.com/gin-gonic/gin"
)

// Thiết lập routes cho Admin (SSR HTML)
func SetupAdminRoutes(router *gin.Engine) {
	dashboardController := adminCtrl.NewDashboardController()

	admin := router.Group("/admin")
	{
		admin.GET("/dashboard", dashboardController.ShowDashboard)
	}
}
