package admin

// ProcessApplicationRequest - Request body cho process application endpoint
type ProcessApplicationRequest struct {
	NewStatus       string `form:"new_status" binding:"required,oneof=Processing Approved Rejected Supplement_Required"`
	Notes           string `form:"notes" binding:"omitempty,max=1000"`
	AssignedStaffID string `form:"assigned_staff_id" binding:"omitempty,uuid"`
}

// CreateServiceRequest - Request body cho create service
type CreateServiceRequest struct {
	Code           string `form:"code" binding:"required"`
	Name           string `form:"name" binding:"required,max=255"`
	Sector         string `form:"sector" binding:"omitempty"`
	DepartmentID   string `form:"department_id" binding:"required,uuid"`
	Description    string `form:"description" binding:"omitempty,max=1000"`
	ProcessingDays int    `form:"processing_days" binding:"omitempty,min=0"`
	Fee            *int   `form:"fee" binding:"omitempty,min=0"`
}

// CreateDepartmentRequest - Request body cho create department
type CreateDepartmentRequest struct {
	Code       string `form:"code" binding:"required"`
	Name       string `form:"name" binding:"required,max=255"`
	LeaderName string `form:"leader_name" binding:"omitempty,max=255"`
	Address    string `form:"address" binding:"omitempty,max=500"`
}

// LoginRequest - Request body cho login
type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}
