package user

// Request query params cho danh sách services (tất cả optional)
type ServiceListRequest struct {
	Page         int    `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	Limit        int    `form:"limit" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
	Keyword      string `form:"keyword" binding:"omitempty" example:"Đăng ký kinh doanh"`
	Sector       string `form:"sector" binding:"omitempty" example:"Health"`
	DepartmentID string `form:"department_id" binding:"omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// Response item cho danh sách services (không có description)
type ServiceListItem struct {
	ID             string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name           string `json:"name" example:"Đăng ký kinh doanh"`
	Code           string `json:"code" example:"SV001"`
	Sector         string `json:"sector" example:"Health"`
	Fee            *int   `json:"fee" example:"500000"`
	ProcessingDays int    `json:"processing_days" example:"5"`
	DepartmentName string `json:"department_name" example:"Sở Y tế"`
}

// Response pagination cho danh sách services
type ServiceListResponse struct {
	Items      []ServiceListItem `json:"items"`
	Page       int               `json:"page" example:"1"`
	Limit      int               `json:"limit" example:"10"`
	TotalItems int64             `json:"total_items" example:"50"`
	TotalPages int               `json:"total_pages" example:"5"`
}

// Response chi tiết service
type ServiceDetailResponse struct {
	ID             string          `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name           string          `json:"name" example:"Đăng ký kinh doanh"`
	Code           string          `json:"code" example:"SV001"`
	Description    string          `json:"description" example:"Dịch vụ đăng ký kinh doanh cho cá nhân và doanh nghiệp"`
	Sector         string          `json:"sector" example:"Health"`
	Fee            *int            `json:"fee" example:"500000"`
	ProcessingDays int             `json:"processing_days" example:"5"`
	Department     *DepartmentInfo `json:"department,omitempty"`
	CreatedAt      string          `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// Thông tin department
type DepartmentInfo struct {
	ID      string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name    string `json:"name" example:"Sở Y tế"`
	Code    string `json:"code" example:"SYT"`
	Address string `json:"address" example:"123 Đường ABC, Quận 1, TP.HCM"`
}
