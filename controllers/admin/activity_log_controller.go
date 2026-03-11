package admin

import (
	"fmt"
	"net/http"
	"strconv"

	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type ActivityLogController struct {
	service *adminSvc.ActivityLogService
}

func NewActivityLogController(service *adminSvc.ActivityLogService) *ActivityLogController {
	return &ActivityLogController{service: service}
}

// GET /admin/activity-logs
func (c *ActivityLogController) List(ctx *gin.Context) {
	action := ctx.Query("action")
	page, error := strconv.Atoi(ctx.Query("page"))
	if error != nil || page <= 0 {
		page = 1
	}

	result, err := c.service.GetList(action, page)
	if err != nil {
		result = nil
	}

	data := utils.GetAdminData(ctx, "Activity Logs", "activity-logs")
	data["Result"] = result
	data["Actions"] = c.service.GetAvailableActions()
	data["Success"] = ctx.Query("success")
	data["Error"] = ctx.Query("error")
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/activity_logs/list.html", data)
}

// POST /admin/activity-logs/cleanup
func (c *ActivityLogController) Cleanup(ctx *gin.Context) {
	days, err := strconv.Atoi(ctx.PostForm("days"))
	if err != nil || days <= 0 {
		days = 30
	}

	deleted, err := c.service.CleanupOldLogs(days)
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/activity-logs")
		return
	}

	setFlashSuccess(ctx, fmt.Sprintf("Deleted %d log(s) older than %d days", deleted, days), "/admin/activity-logs")
}
