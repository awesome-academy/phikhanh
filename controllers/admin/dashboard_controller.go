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
	data["ApplicationCount"] = 0
	data["PendingCount"] = 0
	data["ServiceCount"] = 0
	data["UserCount"] = 0

	utils.RenderHTML(ctx, http.StatusOK, "admin/dashboard.html", data)
}
