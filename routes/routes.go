package routes

import (
	"phikhanh/middlewares"

	"github.com/gin-gonic/gin"
)

// Thiết lập tất cả các routes cho ứng dụng
func SetupRoutes(router *gin.Engine) {
	router.Use(middlewares.CORSMiddleware())

	SetupUserRoutes(router)
	SetupAdminRoutes(router)
}
