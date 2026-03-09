package admin

import (
	"net/http"
	"strconv"
	"time"

	adminDto "phikhanh/dto/admin"
	"phikhanh/models"
	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	service *adminSvc.UserService
}

func NewUserController(service *adminSvc.UserService) *UserController {
	return &UserController{service: service}
}

// GET /admin/users - Danh sách users
func (c *UserController) List(ctx *gin.Context) {
	role := ctx.Query("role")
	page, _ := strconv.Atoi(ctx.Query("page"))

	result, err := c.service.GetList(role, page)
	if err != nil {
		result = nil
	}

	data := utils.GetAdminData(ctx, "Users", "users")
	data["Result"] = result
	data["Roles"] = []string{"all", "staff", "manager", "admin"}
	data["Success"] = ctx.Query("success")
	data["Error"] = ctx.Query("error")
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/users/list.html", data)
}

// GET /admin/users/create - Form tạo user
func (c *UserController) CreateForm(ctx *gin.Context) {
	c.renderForm(ctx, "Create User", &adminDto.UserFormData{}, "/admin/users/create", "Create", "")
}

// POST /admin/users/create - Tạo user mới
func (c *UserController) CreateSave(ctx *gin.Context) {
	var req adminDto.CreateUserRequest

	if err := ctx.ShouldBind(&req); err != nil {
		c.renderFormWithErrors(ctx, "Create User", reqToFormData(req), "/admin/users/create", "Create", utils.FormatValidationErrorsMap(err))
		return
	}

	user := &models.User{
		CitizenID:    req.CitizenID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         models.UserRole(req.Role),
		Phone:        req.Phone,
		Address:      req.Address,
		Gender:       models.Gender(req.Gender),
	}

	if req.DateOfBirth != "" {
		if dob, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			user.DateOfBirth = &dob
		}
	}

	if req.DepartmentID != "" {
		if deptID, err := uuid.Parse(req.DepartmentID); err == nil {
			user.DepartmentID = &deptID
		}
	}

	if err := c.service.Create(user, req.Password); err != nil {
		c.renderFormWithErrors(ctx, "Create User", reqToFormData(req), "/admin/users/create", "Create", map[string]string{"general": formatErrorMessage(err)})
		return
	}

	setFlashSuccess(ctx, "User created successfully", "/admin/users")
}

// GET /admin/users/:id - Chi tiết user
func (c *UserController) ShowDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	detail, err := c.service.GetDetail(id)
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/users")
		return
	}

	data := utils.GetAdminData(ctx, "User Detail", "users")
	data["User"] = detail
	data["Error"] = ctx.Query("error")
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/users/detail.html", data)
}

// GET /admin/users/:id/edit - Form edit user
func (c *UserController) EditForm(ctx *gin.Context) {
	id := ctx.Param("id")

	detail, err := c.service.GetDetail(id)
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/users")
		return
	}

	formData := &adminDto.UserFormData{
		CitizenID:    detail.CitizenID,
		Name:         detail.Name,
		Email:        detail.Email,
		Phone:        detail.Phone,
		Address:      detail.Address,
		DateOfBirth:  detail.DateOfBirth,
		Gender:       detail.Gender,
		Role:         detail.Role,
		DepartmentID: detail.DepartmentID,
	}

	c.renderForm(ctx, "Edit User", formData, "/admin/users/"+id+"/edit", "Update", "")
}

// POST /admin/users/:id/edit - Cập nhật user
func (c *UserController) EditSave(ctx *gin.Context) {
	id := ctx.Param("id")

	var req adminDto.UpdateUserRequest

	if err := ctx.ShouldBind(&req); err != nil {
		c.renderFormWithErrors(ctx, "Edit User", updateReqToFormData(req), "/admin/users/"+id+"/edit", "Update", utils.FormatValidationErrorsMap(err))
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		setFlashError(ctx, "Invalid user ID", "/admin/users")
		return
	}

	userRole, _ := ctx.Get("admin_role")
	roleStr, _ := userRole.(string)

	user := &models.User{
		BaseModel: models.BaseModel{ID: userID},
		CitizenID: req.CitizenID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Gender:    models.Gender(req.Gender),
		Role:      models.UserRole(req.Role),
	}

	if req.DateOfBirth != "" {
		if dob, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			user.DateOfBirth = &dob
		}
	}

	if req.DepartmentID != "" {
		if deptID, err := uuid.Parse(req.DepartmentID); err == nil {
			user.DepartmentID = &deptID
		}
	}

	if err := c.service.ValidateAndPrepareForUpdate(user, roleStr); err != nil {
		c.renderFormWithErrors(ctx, "Edit User", updateReqToFormData(req), "/admin/users/"+id+"/edit", "Update", map[string]string{"general": formatErrorMessage(err)})
		return
	}

	if err := c.service.Update(user); err != nil {
		c.renderFormWithErrors(ctx, "Edit User", updateReqToFormData(req), "/admin/users/"+id+"/edit", "Update", map[string]string{"general": formatErrorMessage(err)})
		return
	}

	setFlashSuccess(ctx, "User updated successfully", "/admin/users")
}

// POST /admin/users/:id/delete - Delete user
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/users/"+id)
		return
	}

	setFlashSuccess(ctx, "User deleted successfully", "/admin/users")
}

// Helper functions
// renderFormWithErrors - Render form với map errors cho từng field
func (c *UserController) renderFormWithErrors(ctx *gin.Context, title string, formData *adminDto.UserFormData, action, label string, fieldErrors map[string]string) {
	userRole, _ := ctx.Get("admin_role")
	roleStr, _ := userRole.(string)

	roles := c.service.GetAvailableRoles()
	genders := c.service.GetAvailableGenders()
	departments, _ := c.service.GetDepartments()

	data := utils.GetAdminData(ctx, title, "users")
	data["User"] = formData
	data["Roles"] = roles
	data["Genders"] = genders
	data["Departments"] = departments
	data["FormAction"] = action
	data["SubmitLabel"] = label
	data["FieldErrors"] = fieldErrors
	data["UserRole"] = roleStr
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/users/form.html", data)
}

func (c *UserController) renderForm(ctx *gin.Context, title string, formData *adminDto.UserFormData, action, label, errMsg string) {
	errMap := map[string]string{}
	if errMsg != "" {
		errMap["general"] = errMsg
	}
	c.renderFormWithErrors(ctx, title, formData, action, label, errMap)
}

// Helper: convert CreateUserRequest -> UserFormData
func reqToFormData(req adminDto.CreateUserRequest) *adminDto.UserFormData {
	return &adminDto.UserFormData{
		CitizenID:    req.CitizenID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		Role:         req.Role,
		DepartmentID: req.DepartmentID,
	}
}

// Helper: convert UpdateUserRequest -> UserFormData
func updateReqToFormData(req adminDto.UpdateUserRequest) *adminDto.UserFormData {
	return &adminDto.UserFormData{
		CitizenID:    req.CitizenID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		Role:         req.Role,
		DepartmentID: req.DepartmentID,
	}
}
