package admin

import (
	"net/http"
	"phikhanh/models"
	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceController struct {
	service *adminSvc.ServiceAdminService
}

func NewServiceController(service *adminSvc.ServiceAdminService) *ServiceController {
	return &ServiceController{service: service}
}

// List - Danh sách services
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

// Detail - Chi tiết service
func (c *ServiceController) Detail(ctx *gin.Context) {
	service, ok := c.findService(ctx)
	if !ok {
		return
	}

	data := utils.GetAdminData(ctx, "Service Detail", "services")
	data["Service"] = service
	data["CreatedAt"] = service.CreatedAt.Format(time.RFC3339)
	data["UpdatedAt"] = service.UpdatedAt.Format(time.RFC3339)

	utils.RenderHTML(ctx, http.StatusOK, "admin/services/detail.html", data)
}

// CreateForm - Hiển thị form tạo service
func (c *ServiceController) CreateForm(ctx *gin.Context) {
	utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html", c.formData(ctx, "Add New Service", &models.Service{}, "/admin/services/create", "Create Service", ""))
}

// CreateSave - Xử lý tạo service
func (c *ServiceController) CreateSave(ctx *gin.Context) {
	service, err := c.service.BindForm(ctx)
	if err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html", c.formData(ctx, "Add New Service", service, "/admin/services/create", "Create Service", err.Error()))
		return
	}

	if err := c.service.Create(service); err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html", c.formData(ctx, "Add New Service", service, "/admin/services/create", "Create Service", err.Error()))
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/services?success=Service+created+successfully")
}

// EditForm - Hiển thị form edit service
func (c *ServiceController) EditForm(ctx *gin.Context) {
	service, ok := c.findService(ctx)
	if !ok {
		return
	}

	utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
		c.formData(ctx, "Edit Service", service, "/admin/services/"+service.ID.String()+"/edit", "Save Changes", ""))
}

// EditSave - Xử lý update service
func (c *ServiceController) EditSave(ctx *gin.Context) {
	existing, ok := c.findService(ctx)
	if !ok {
		return
	}

	updated, err := c.service.BindForm(ctx)
	if err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
			c.formData(ctx, "Edit Service", updated, "/admin/services/"+existing.ID.String()+"/edit", "Save Changes", err.Error()))
		return
	}

	updated.ID = existing.ID
	if err := c.service.Update(updated); err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/services/form.html",
			c.formData(ctx, "Edit Service", updated, "/admin/services/"+existing.ID.String()+"/edit", "Save Changes", err.Error()))
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/services?success=Service+updated+successfully")
}

// Delete - Xóa service
func (c *ServiceController) Delete(ctx *gin.Context) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return
	}

	if err := c.service.Delete(id); err != nil {
		ctx.Redirect(http.StatusFound, "/admin/services?error=Failed+to+delete+service")
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/services?success=Service+deleted+successfully")
}

// --- Private helpers ---

func (c *ServiceController) findService(ctx *gin.Context) (*models.Service, bool) {
	id, ok := parseServiceID(ctx)
	if !ok {
		return nil, false
	}

	service, err := c.service.GetByID(id)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/admin/services?error=Service+not+found")
		return nil, false
	}

	return service, true
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

func parseServiceID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Redirect(http.StatusFound, "/admin/services?error=Invalid+service+ID")
		return uuid.Nil, false
	}
	return id, true
}
