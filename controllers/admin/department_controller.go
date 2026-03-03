package admin

import (
	"net/http"
	"phikhanh/models"
	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct {
	service *adminSvc.DepartmentService
}

func NewDepartmentController(service *adminSvc.DepartmentService) *DepartmentController {
	return &DepartmentController{service: service}
}

// GET /admin/departments
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

// GET /admin/departments/:id
func (c *DepartmentController) Detail(ctx *gin.Context) {
	id, ok := parseDepartmentID(ctx)
	if !ok {
		return
	}

	detail, err := c.service.GetDetail(id)
	if err != nil {
		setFlashError(ctx, "Department not found", redirectDepartments)
		return
	}

	data := utils.GetAdminData(ctx, "Department Detail", "departments")
	data["Department"] = detail

	utils.RenderHTML(ctx, http.StatusOK, "admin/departments/detail.html", data)
}

// GET /admin/departments/create
func (c *DepartmentController) CreateForm(ctx *gin.Context) {
	c.renderForm(ctx, "Add New Department", &models.Department{}, "/admin/departments/create", "Create Department", "")
}

// POST /admin/departments/create
func (c *DepartmentController) CreateSave(ctx *gin.Context) {
	department, err := c.service.BindForm(ctx)
	if err != nil {
		c.renderForm(ctx, "Add New Department", department, "/admin/departments/create", "Create Department", formatErrorMessage(err))
		return
	}

	if err := c.service.Create(department); err != nil {
		c.renderForm(ctx, "Add New Department", department, "/admin/departments/create", "Create Department", formatErrorMessage(err))
		return
	}

	setFlashSuccess(ctx, "Department created successfully", redirectDepartments)
}

// GET /admin/departments/:id/edit
func (c *DepartmentController) EditForm(ctx *gin.Context) {
	id, ok := parseDepartmentID(ctx)
	if !ok {
		return
	}

	department, err := c.service.GetByID(id)
	if err != nil {
		setFlashError(ctx, "Department not found", redirectDepartments)
		return
	}

	c.renderForm(ctx, "Edit Department", department, "/admin/departments/"+department.ID.String()+"/edit", "Save Changes", "")
}

// POST /admin/departments/:id/edit
func (c *DepartmentController) EditSave(ctx *gin.Context) {
	id, ok := parseDepartmentID(ctx)
	if !ok {
		return
	}

	updated, err := c.service.BindForm(ctx)
	if err != nil {
		c.renderForm(ctx, "Edit Department", updated, "/admin/departments/"+id.String()+"/edit", "Save Changes", formatErrorMessage(err))
		return
	}

	updated.ID = id
	if err := c.service.Update(updated); err != nil {
		c.renderForm(ctx, "Edit Department", updated, "/admin/departments/"+id.String()+"/edit", "Save Changes", formatErrorMessage(err))
		return
	}

	setFlashSuccess(ctx, "Department updated successfully", redirectDepartments)
}

// POST /admin/departments/:id/delete
func (c *DepartmentController) Delete(ctx *gin.Context) {
	id, ok := parseDepartmentID(ctx)
	if !ok {
		return
	}

	if err := c.service.Delete(id); err != nil {
		setFlashError(ctx, "Failed to delete department", redirectDepartments)
		return
	}

	setFlashSuccess(ctx, "Department deleted successfully", redirectDepartments)
}

// Helper function để render form
func (c *DepartmentController) renderForm(ctx *gin.Context, title string, department *models.Department, action, label, errMsg string) {
	utils.RenderHTML(ctx, http.StatusOK, "admin/departments/form.html", c.formData(ctx, title, department, action, label, errMsg))
}

func (c *DepartmentController) formData(ctx *gin.Context, title string, department *models.Department, action, label, errMsg string) gin.H {
	data := utils.GetAdminData(ctx, title, "departments")
	data["Department"] = department
	data["FormAction"] = action
	data["SubmitLabel"] = label
	data["Error"] = errMsg
	data["Success"] = ctx.Query("success")
	data["CsrfToken"] = getCsrfToken(ctx)
	return data
}
