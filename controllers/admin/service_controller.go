package admin

import (
	"fmt"
	"net/http"

	"phikhanh/models"
	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	service        *adminSvc.ServiceAdminService
	activityLogSvc *adminSvc.ActivityLogService
}

func NewServiceController(service *adminSvc.ServiceAdminService, activityLogSvc *adminSvc.ActivityLogService) *ServiceController {
	return &ServiceController{service: service, activityLogSvc: activityLogSvc}
}

// GET /admin/services
func (c *ServiceController) List(ctx *gin.Context) {
	services, err := c.service.GetAll()
	if err != nil {
		services = []models.Service{}
	}

	data := utils.GetAdminData(ctx, "Services", "services")
	data["Services"] = services
	data["Success"] = ctx.Query("success")
	data["Error"] = ctx.Query("error")

	utils.RenderHTML(ctx, http.StatusOK, "admin/services/list.html", data)
}

// GET /admin/services/:id - Hiển thị chi tiết service
func (c *ServiceController) Detail(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	detail, err := c.service.GetDetail(id)
	if err != nil {
		setFlashError(ctx, "Service not found", redirectServices)
		return
	}

	data := utils.GetAdminData(ctx, "Service Detail", "services")
	data["Service"] = detail
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/services/detail.html", data)
}

// GET /admin/services/create
func (c *ServiceController) CreateForm(ctx *gin.Context) {
	c.renderForm(ctx, "Add New Service", &models.Service{}, "/admin/services/create", "Create Service", "")
}

// POST /admin/services/create
func (c *ServiceController) CreateSave(ctx *gin.Context) {
	service, err := c.service.BindForm(ctx)
	if err != nil {
		c.renderForm(ctx, "Add New Service", service, "/admin/services/create", "Create Service", formatErrorMessage(err))
		return
	}
	if err := c.service.Create(service); err != nil {
		c.renderForm(ctx, "Add New Service", service, "/admin/services/create", "Create Service", formatErrorMessage(err))
		return
	}

	actorID, _ := utils.ExtractAdminID(ctx)
	c.activityLogSvc.RecordActivity(
		actorID.String(),
		models.ActionCreateService,
		service.ID.String(),
		fmt.Sprintf("Created service: %s (%s)", service.Name, service.Code),
		ctx.ClientIP(),
	)

	setFlashSuccess(ctx, "Service created successfully", redirectServices)
}

// GET /admin/services/:id/edit
func (c *ServiceController) EditForm(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	service, err := c.service.GetByID(id)
	if err != nil {
		setFlashError(ctx, "Service not found", redirectServices)
		return
	}

	c.renderForm(ctx, "Edit Service", service, "/admin/services/"+service.ID.String()+"/edit", "Save Changes", "")
}

// POST /admin/services/:id/edit
func (c *ServiceController) EditSave(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}
	updated, err := c.service.BindForm(ctx)
	if err != nil {
		c.renderForm(ctx, "Edit Service", updated, "/admin/services/"+id.String()+"/edit", "Save Changes", formatErrorMessage(err))
		return
	}
	updated.ID = id
	if err := c.service.Update(updated); err != nil {
		c.renderForm(ctx, "Edit Service", updated, "/admin/services/"+id.String()+"/edit", "Save Changes", formatErrorMessage(err))
		return
	}

	actorID, _ := utils.ExtractAdminID(ctx)
	c.activityLogSvc.RecordActivity(
		actorID.String(),
		models.ActionUpdateService,
		id.String(),
		fmt.Sprintf("Updated service: %s (%s)", updated.Name, updated.Code),
		ctx.ClientIP(),
	)

	setFlashSuccess(ctx, "Service updated successfully", redirectServices)
}

// POST /admin/services/:id/delete
func (c *ServiceController) Delete(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	svc, _ := c.service.GetByID(id)

	if err := c.service.Delete(id); err != nil {
		setFlashError(ctx, "Failed to delete service", redirectServices)
		return
	}

	actorID, _ := utils.ExtractAdminID(ctx)
	desc := "Deleted service ID: " + id.String()
	if svc != nil {
		desc = fmt.Sprintf("Deleted service: %s (%s)", svc.Name, svc.Code)
	}
	c.activityLogSvc.RecordActivity(actorID.String(), models.ActionDeleteService, id.String(), desc, ctx.ClientIP())

	setFlashSuccess(ctx, "Service deleted successfully", redirectServices)
}

// Helper function để render form
func (c *ServiceController) renderForm(ctx *gin.Context, title string, service *models.Service, action, label, errMsg string) {
	utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html", c.formData(ctx, title, service, action, label, errMsg))
}

func (c *ServiceController) formData(ctx *gin.Context, title string, service *models.Service, action, label, errMsg string) gin.H {
	departments, _ := c.service.GetDepartments()
	data := utils.GetAdminData(ctx, title, "services")
	data["Service"] = service
	data["Departments"] = departments
	data["Sectors"] = adminSvc.AvailableSectors
	data["FormAction"] = action
	data["SubmitLabel"] = label
	data["Error"] = errMsg
	data["Success"] = ctx.Query("success")
	data["CsrfToken"] = getCsrfToken(ctx)
	return data
}
