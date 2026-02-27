package admin

import (
	"net/http"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (c *DashboardController) ShowDashboard(ctx *gin.Context) {
	data := utils.GetAdminData(ctx, "Dashboard", "dashboard")
	data["ApplicationCount"] = 24
	data["PendingCount"] = 5
	data["ServiceCount"] = 12
	data["UserCount"] = 156

	ctx.HTML(http.StatusOK, "admin/dashboard.html", data)
}
