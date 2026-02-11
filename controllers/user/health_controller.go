package user

import (
	"net/http"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// CheckHealth - Kiểm tra trạng thái hệ thống
// GET /api/v1/health
func (c *HealthController) CheckHealth(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  "healthy",
		"message": "Public Service Management API is running",
		"version": "1.0.0",
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Health check successful", response)
}
