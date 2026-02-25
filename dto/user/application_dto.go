package user

// AttachmentDTO - Thông tin file đính kèm
type AttachmentDTO struct {
	FilePath string `json:"file_path" binding:"required" example:"/assets/images/20240115_a1b2c3d4.jpg"`
	FileName string `json:"file_name" binding:"required" example:"document.jpg"`
}

// SubmitAppRequest - Request nộp hồ sơ
type SubmitAppRequest struct {
	ServiceID   string          `json:"service_id" binding:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Attachments []AttachmentDTO `json:"attachments" binding:"omitempty,dive,required"`
}

// SubmitAppResponse - Response sau khi nộp hồ sơ
type SubmitAppResponse struct {
	ID   string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Code string `json:"code" example:"HS-20240115-a1b2c3d4"`
}

// MyAppListRequest - Query params cho danh sách hồ sơ của tôi
type MyAppListRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
	Status string `form:"status" binding:"omitempty,oneof=Received Processing Supplement_Required Approved Rejected" example:"Received"`
}

// MyAppItemResponse - Response item cho danh sách hồ sơ
type MyAppItemResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Code        string `json:"code" example:"HS-20240115-a1b2c3d4"`
	ServiceName string `json:"service_name" example:"Đăng ký kinh doanh"`
	Status      string `json:"status" example:"Received"`
	CreatedAt   string `json:"created_at" example:"2024-01-15T10:00:00Z"`
}

// MyAppListResponse - Paginated response
type MyAppListResponse struct {
	Items      []MyAppItemResponse `json:"items"`
	Page       int                 `json:"page" example:"1"`
	Limit      int                 `json:"limit" example:"10"`
	TotalItems int64               `json:"total_items" example:"20"`
	TotalPages int                 `json:"total_pages" example:"2"`
}
