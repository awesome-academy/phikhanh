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

// GET /admin/services/:id
func (c *ServiceController) Detail(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	detail, err := c.service.GetDetail(id)
	if err != nil {
		ctx.Redirect(http.StatusFound, redirectServices+"?error=Service+not+found")
		return
	}

	data := utils.GetAdminData(ctx, "Service Detail", "services")
	data["Service"] = detail

	utils.RenderHTML(ctx, http.StatusOK, "admin/services/detail.html", data)
}

// GET /admin/services/create
func (c *ServiceController) CreateForm(ctx *gin.Context) {
	utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
		c.formData(ctx, "Add New Service", &models.Service{}, "/admin/services/create", "Create Service", ""))
}

// POST /admin/services/create
func (c *ServiceController) CreateSave(ctx *gin.Context) {
	service, err := c.service.BindForm(ctx)
	if err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
			c.formData(ctx, "Add New Service", service, "/admin/services/create", "Create Service", err.Error()))
		return
	}

	if err := c.service.Create(service); err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
			c.formData(ctx, "Add New Service", service, "/admin/services/create", "Create Service", err.Error()))
		return
	}

	ctx.Redirect(http.StatusFound, redirectServices+"?success=Service+created+successfully")
}

// GET /admin/services/:id/edit
func (c *ServiceController) EditForm(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	service, err := c.service.GetByID(id)
	if err != nil {
		ctx.Redirect(http.StatusFound, redirectServices+"?error=Service+not+found")
		return
	}

	utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
		c.formData(ctx, "Edit Service", service, "/admin/services/"+service.ID.String()+"/edit", "Save Changes", ""))
}

// POST /admin/services/:id/edit
func (c *ServiceController) EditSave(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	updated, err := c.service.BindForm(ctx)
	if err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
			c.formData(ctx, "Edit Service", updated, "/admin/services/"+id.String()+"/edit", "Save Changes", err.Error()))
		return
	}

	updated.ID = id
	if err := c.service.Update(updated); err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
			c.formData(ctx, "Edit Service", updated, "/admin/services/"+id.String()+"/edit", "Save Changes", err.Error()))
		return
	}

	ctx.Redirect(http.StatusFound, redirectServices+"?success=Service+updated+successfully")
}

// POST /admin/services/:id/delete
func (c *ServiceController) Delete(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	if err := c.service.Delete(id); err != nil {
		ctx.Redirect(http.StatusFound, redirectServices+"?error=Failed+to+delete+service")
		return
	}

	ctx.Redirect(http.StatusFound, redirectServices+"?success=Service+deleted+successfully")
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
	return data
}
