package admin

import (
	"fmt"
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
	service        *adminSvc.UserService
	activityLogSvc *adminSvc.ActivityLogService
}

func NewUserController(service *adminSvc.UserService, activityLogSvc *adminSvc.ActivityLogService) *UserController {
	return &UserController{service: service, activityLogSvc: activityLogSvc}
}

// GET /admin/users - Danh sách users
func (c *UserController) List(ctx *gin.Context) {
	role := ctx.Query("role")
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	result, err := c.service.GetList(role, page)
	if err != nil {
		result = nil
	}

	data := utils.GetAdminData(ctx, "Users", "users")
	data["Result"] = result
	data["Roles"] = []string{"all", "citizen", "staff", "manager", "admin"}
	data["Success"] = ctx.Query("success")
	data["Error"] = ctx.Query("error")
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/users/list.html", data)
}

// GET /admin/users/create - Form tạo user
func (c *UserController) CreateForm(ctx *gin.Context) {
	c.renderForm(ctx, "Create User", &adminDto.UserFormData{}, "/admin/users/create", "Create", true, "")
}

// POST /admin/users/create - Tạo user mới
func (c *UserController) CreateSave(ctx *gin.Context) {
	var req adminDto.CreateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		c.renderFormWithErrors(ctx, "Create User", reqToFormData(req), "/admin/users/create", "Create", true, utils.FormatValidationErrorsMap(err))
		return
	}

	user := &models.User{
		CitizenID: req.CitizenID, Name: req.Name, Email: req.Email,
		PasswordHash: req.Password, Role: models.UserRole(req.Role),
		Phone: req.Phone, Address: req.Address, Gender: models.Gender(req.Gender),
	}
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			c.renderFormWithErrors(ctx, "Create User", reqToFormData(req), "/admin/users/create", "Create", true,
				map[string]string{"date_of_birth": "Invalid date format, expected YYYY-MM-DD"})
			return
		}
		user.DateOfBirth = &dob
	}
	if req.DepartmentID != "" {
		deptID, err := uuid.Parse(req.DepartmentID)
		if err != nil {
			c.renderFormWithErrors(ctx, "Create User", reqToFormData(req), "/admin/users/create", "Create", true,
				map[string]string{"department_id": "Invalid department ID"})
			return
		}
		user.DepartmentID = &deptID
	}

	if err := c.service.Create(user, req.Password); err != nil {
		c.renderFormWithErrors(ctx, "Create User", reqToFormData(req), "/admin/users/create", "Create", true, map[string]string{"general": formatErrorMessage(err)})
		return
	}

	// Record activity
	actorID, _ := utils.ExtractAdminID(ctx)
	c.activityLogSvc.RecordActivity(
		actorID.String(),
		models.ActionCreateUser,
		user.ID.String(),
		fmt.Sprintf("Created user: %s (%s) with role %s", user.Name, user.CitizenID, user.Role),
		ctx.ClientIP(),
	)

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

	c.renderForm(ctx, "Edit User", formData, "/admin/users/"+id+"/edit", "Update", false, "")
}

// POST /admin/users/:id/edit - Cập nhật user
func (c *UserController) EditSave(ctx *gin.Context) {
	id := ctx.Param("id")
	var req adminDto.UpdateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		c.renderFormWithErrors(ctx, "Edit User", updateReqToFormData(req), "/admin/users/"+id+"/edit", "Update", false, utils.FormatValidationErrorsMap(err))
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
		CitizenID: req.CitizenID, Name: req.Name, Email: req.Email,
		Phone: req.Phone, Address: req.Address,
		Gender: models.Gender(req.Gender), Role: models.UserRole(req.Role),
	}
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			c.renderFormWithErrors(ctx, "Edit User", updateReqToFormData(req), "/admin/users/"+id+"/edit", "Update", false,
				map[string]string{"date_of_birth": "Invalid date format, expected YYYY-MM-DD"})
			return
		}
		user.DateOfBirth = &dob
	}
	if req.DepartmentID != "" {
		deptID, err := uuid.Parse(req.DepartmentID)
		if err != nil {
			c.renderFormWithErrors(ctx, "Edit User", updateReqToFormData(req), "/admin/users/"+id+"/edit", "Update", false,
				map[string]string{"department_id": "Invalid department ID"})
			return
		}
		user.DepartmentID = &deptID
	}

	if err := c.service.ValidateAndPrepareForUpdate(user, roleStr); err != nil {
		c.renderFormWithErrors(ctx, "Edit User", updateReqToFormData(req), "/admin/users/"+id+"/edit", "Update", false, map[string]string{"general": formatErrorMessage(err)})
		return
	}
	if err := c.service.Update(user); err != nil {
		c.renderFormWithErrors(ctx, "Edit User", updateReqToFormData(req), "/admin/users/"+id+"/edit", "Update", false, map[string]string{"general": formatErrorMessage(err)})
		return
	}

	// Record activity
	actorID, _ := utils.ExtractAdminID(ctx)
	c.activityLogSvc.RecordActivity(
		actorID.String(),
		models.ActionUpdateUser,
		id,
		fmt.Sprintf("Updated user: %s (%s)", req.Name, req.CitizenID),
		ctx.ClientIP(),
	)

	setFlashSuccess(ctx, "User updated successfully", "/admin/users")
}

// POST /admin/users/:id/delete - Delete user
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get user info before delete for logging
	detail, _ := c.service.GetDetail(id)

	if err := c.service.Delete(id); err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/users/"+id)
		return
	}

	// Record activity
	actorID, _ := utils.ExtractAdminID(ctx)
	desc := "Deleted user ID: " + id
	if detail != nil {
		desc = fmt.Sprintf("Deleted user: %s (%s)", detail.Name, detail.CitizenID)
	}
	c.activityLogSvc.RecordActivity(actorID.String(), models.ActionDeleteUser, id, desc, ctx.ClientIP())

	setFlashSuccess(ctx, "User deleted successfully", "/admin/users")
}

// Helper functions
// renderFormWithErrors - Render form với map errors cho từng field
func (c *UserController) renderFormWithErrors(ctx *gin.Context, title string, formData *adminDto.UserFormData, action, label string, isCreate bool, fieldErrors map[string]string) {
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
	data["IsCreate"] = isCreate
	data["FieldErrors"] = fieldErrors
	data["UserRole"] = roleStr
	data["CsrfToken"] = getCsrfToken(ctx)

	utils.RenderHTML(ctx, http.StatusOK, "admin/users/form.html", data)
}

func (c *UserController) renderForm(ctx *gin.Context, title string, formData *adminDto.UserFormData, action, label string, isCreate bool, errMsg string) {
	errMap := map[string]string{}
	if errMsg != "" {
		errMap["general"] = errMsg
	}
	c.renderFormWithErrors(ctx, title, formData, action, label, isCreate, errMap)
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
