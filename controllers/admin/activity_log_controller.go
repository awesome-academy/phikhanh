package admin

import (
	"net/http"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type ActivityLogController struct{}

func NewActivityLogController() *ActivityLogController {
	return &ActivityLogController{}
}

func (c *ActivityLogController) ShowActivityLogs(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/activity_logs.html", utils.GetAdminData(ctx, "Activity Log", "activity-logs"))
}
