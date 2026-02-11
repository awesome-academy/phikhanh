package routes

import (
	userCtrl "phikhanh/controllers/user"

	"github.com/gin-gonic/gin"
)

// Thiết lập routes cho User API (JSON)
func SetupUserRoutes(router *gin.Engine) {
	healthController := userCtrl.NewHealthController()

	api := router.Group("/api/v1")
	{
		api.GET("/health", healthController.CheckHealth)
	}
}
