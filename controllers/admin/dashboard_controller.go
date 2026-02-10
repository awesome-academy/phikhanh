package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// ShowDashboard - Hiển thị trang dashboard admin
// GET /admin/dashboard
func (c *DashboardController) ShowDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "dashboard.html", gin.H{})
}
