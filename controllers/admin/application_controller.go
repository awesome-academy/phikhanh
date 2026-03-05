package admin

import (
	"net/http"
	"strconv"

	adminDto "phikhanh/dto/admin"
	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type ApplicationController struct {
	service *adminSvc.ApplicationAdminService
}

func NewApplicationController(service *adminSvc.ApplicationAdminService) *ApplicationController {
	return &ApplicationController{service: service}
}

// GET /admin/applications - Danh sách applications
// Query params: status (filter), page (pagination)
// Staff chỉ thấy applications assigned to themselves (auto filter)
func (c *ApplicationController) List(ctx *gin.Context) {
	status := ctx.Query("status")
	page, _ := strconv.Atoi(ctx.Query("page"))

	// Nếu user là staff hoặc request có assigned_to_me=true, filter theo logged-in user
	// Admin & Manager thấy tất cả applications
	var assignedToUserID *string
	userRole, _ := ctx.Get("admin_role")
	roleStr, _ := userRole.(string)

	if roleStr == "staff" || ctx.Query("assigned_to_me") == "true" {
		adminID, svcErr := utils.ExtractAdminID(ctx)
		if svcErr == nil {
			idStr := adminID.String()
			assignedToUserID = &idStr
		}
	}

	result, err := c.service.GetList(status, assignedToUserID, page)
	if err != nil {
		result = nil
	}

	data := utils.GetAdminData(ctx, "Applications", "applications")
	data["Result"] = result
	data["Statuses"] = []string{
		"Received", "Processing", "Supplement_Required", "Approved", "Rejected",
	}
	data["UserRole"] = roleStr
	data["Success"] = ctx.Query("success")
	data["Error"] = ctx.Query("error")
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/applications/list.html", data)
}

// GET /admin/applications/:id - Chi tiết application
// Hiển thị: applicant info, service info, attachments, history timeline, process form
func (c *ApplicationController) ShowDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	detail, err := c.service.GetDetail(id)
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/applications")
		return
	}

	// Lấy danh sách staff có thể assign
	availableStaff, _ := c.service.GetAvailableStaff()

	userRole, _ := ctx.Get("admin_role")
	roleStr, _ := userRole.(string)

	// Lấy danh sách status tiếp theo (state machine based)
	nextStatuses := c.service.GetNextStatuses(detail.Status)

	data := utils.GetAdminData(ctx, "Application Detail", "applications")
	data["Application"] = detail
	data["Statuses"] = nextStatuses
	data["AvailableStaff"] = availableStaff
	data["UserRole"] = roleStr
	data["Error"] = ctx.Query("error")
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/applications/detail.html", data)
}

// POST /admin/applications/:id/process - Xử lý application
func (c *ApplicationController) Process(ctx *gin.Context) {
	appID := ctx.Param("id")

	adminID, svcErr := utils.ExtractAdminID(ctx)
	if svcErr != nil {
		setFlashError(ctx, "Unauthorized", "/admin/applications")
		return
	}

	var req adminDto.ProcessApplicationRequest

	if err := ctx.ShouldBind(&req); err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/applications/"+appID)
		return
	}

	userRole, _ := ctx.Get("admin_role")
	roleStr, _ := userRole.(string)

	var assignedStaffID *string
	if roleStr == "admin" || roleStr == "manager" {
		if req.AssignedStaffID != "" {
			assignedStaffID = &req.AssignedStaffID
		}
	}

	if err := c.service.ProcessApplication(appID, req.NewStatus, assignedStaffID, req.Notes, adminID.String()); err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			ctx.Redirect(http.StatusFound, "/admin/applications/"+appID+"?error="+svcErr.Message)
			return
		}
		ctx.Redirect(http.StatusFound, "/admin/applications/"+appID+"?error=Failed+to+process+application")
		return
	}

	setFlashSuccess(ctx, "Application processed successfully", "/admin/applications")
}
