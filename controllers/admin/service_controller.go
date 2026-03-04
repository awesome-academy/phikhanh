package admin

import (
	"net/http"
	"phikhanh/models"
	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	service *adminSvc.ServiceAdminService
}

func NewServiceController(service *adminSvc.ServiceAdminService) *ServiceController {
	return &ServiceController{service: service}
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

	setFlashSuccess(ctx, "Service updated successfully", redirectServices)
}

// POST /admin/services/:id/delete
func (c *ServiceController) Delete(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	if err := c.service.Delete(id); err != nil {
		setFlashError(ctx, "Failed to delete service", redirectServices)
		return
	}

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
