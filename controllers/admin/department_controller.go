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

type DepartmentController struct {
	service *adminSvc.DepartmentService
}

func NewDepartmentController(service *adminSvc.DepartmentService) *DepartmentController {
	return &DepartmentController{service: service}
}

func (c *DepartmentController) List(ctx *gin.Context) {
	departments, err := c.service.GetAll()
	if err != nil {
		departments = []models.Department{}
	}

	data := utils.GetAdminData(ctx, "Departments", "departments")
	data["Departments"] = departments
	data["Success"] = ctx.Query("success")
	data["Error"] = ctx.Query("error")

	utils.RenderHTML(ctx, http.StatusOK, "admin/departments/list.html", data)
}

func (c *DepartmentController) Detail(ctx *gin.Context) {
	department, ok := c.findDepartment(ctx)
	if !ok {
		return
	}

	data := utils.GetAdminData(ctx, "Department Detail", "departments")
	data["Department"] = department
	data["CreatedAt"] = department.CreatedAt.Format(time.RFC3339)
	data["UpdatedAt"] = department.UpdatedAt.Format(time.RFC3339)

	utils.RenderHTML(ctx, http.StatusOK, "admin/departments/detail.html", data)
}

func (c *DepartmentController) CreateForm(ctx *gin.Context) {
	utils.RenderHTML(ctx, http.StatusOK, "admin/departments/form.html",
		c.formData(ctx, "Add New Department", &models.Department{}, "/admin/departments/create", "Create Department", ""))
}

func (c *DepartmentController) CreateSave(ctx *gin.Context) {
	department, err := c.service.BindForm(ctx)
	if err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/departments/form.html",
			c.formData(ctx, "Add New Department", department, "/admin/departments/create", "Create Department", err.Error()))
		return
	}

	if err := c.service.Create(department); err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/departments/form.html",
			c.formData(ctx, "Add New Department", department, "/admin/departments/create", "Create Department", err.Error()))
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/departments?success=Department+created+successfully")
}

func (c *DepartmentController) EditForm(ctx *gin.Context) {
	department, ok := c.findDepartment(ctx)
	if !ok {
		return
	}

	utils.RenderHTML(ctx, http.StatusOK, "admin/departments/form.html",
		c.formData(ctx, "Edit Department", department, "/admin/departments/"+department.ID.String()+"/edit", "Save Changes", ""))
}

func (c *DepartmentController) EditSave(ctx *gin.Context) {
	existing, ok := c.findDepartment(ctx)
	if !ok {
		return
	}

	updated, err := c.service.BindForm(ctx)
	if err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/departments/form.html",
			c.formData(ctx, "Edit Department", updated, "/admin/departments/"+existing.ID.String()+"/edit", "Save Changes", err.Error()))
		return
	}

	updated.ID = existing.ID
	if err := c.service.Update(updated); err != nil {
		utils.RenderHTML(ctx, http.StatusOK, "admin/departments/form.html",
			c.formData(ctx, "Edit Department", updated, "/admin/departments/"+existing.ID.String()+"/edit", "Save Changes", err.Error()))
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/departments?success=Department+updated+successfully")
}

func (c *DepartmentController) Delete(ctx *gin.Context) {
	id, ok := parseDepartmentID(ctx)
	if !ok {
		return
	}

	if err := c.service.Delete(id); err != nil {
		ctx.Redirect(http.StatusFound, "/admin/departments?error=Failed+to+delete+department")
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/departments?success=Department+deleted+successfully")
}

// --- Private helpers ---

func (c *DepartmentController) findDepartment(ctx *gin.Context) (*models.Department, bool) {
	id, ok := parseDepartmentID(ctx)
	if !ok {
		return nil, false
	}

	department, err := c.service.GetByID(id)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/admin/departments?error=Department+not+found")
		return nil, false
	}

	return department, true
}

func (c *DepartmentController) formData(ctx *gin.Context, title string, department *models.Department, action, label, errMsg string) gin.H {
	data := utils.GetAdminData(ctx, title, "departments")
	data["Department"] = department
	data["FormAction"] = action
	data["SubmitLabel"] = label
	data["Error"] = errMsg
	return data
}

func parseDepartmentID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Redirect(http.StatusFound, "/admin/departments?error=Invalid+department+ID")
		return uuid.Nil, false
	}
	return id, true
}
