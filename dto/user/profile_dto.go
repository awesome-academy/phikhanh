package user

// Response thông tin hồ sơ công dân
type ProfileResponse struct {
	ID            string  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CitizenID     string  `json:"citizen_id" example:"001234567890"`
	Name          string  `json:"name" example:"Nguyễn Văn A"`
	Email         string  `json:"email" example:"nguyenvana@example.com"`
	Phone         string  `json:"phone" example:"0901234567"`
	Address       string  `json:"address" example:"123 Đường ABC, Quận 1, TP.HCM"`
	DateOfBirth   *string `json:"date_of_birth" example:"1990-01-15"`
	Gender        string  `json:"gender" example:"male"`
	Role          string  `json:"role" example:"citizen"`
	IsEmailNotify bool    `json:"is_email_notify" example:"true"`
	CreatedAt     string  `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// Request cập nhật profile
type UpdateProfileRequest struct {
	Name          string  `json:"name" binding:"required" example:"Nguyễn Văn A"`
	Phone         string  `json:"phone" binding:"required,vn_phone" example:"0901234567"`
	Address       string  `json:"address" example:"123 Đường ABC, Quận 1, TP.HCM"`
	DateOfBirth   *string `json:"date_of_birth" binding:"omitempty,past_date" example:"1990-01-15"`
	Gender        string  `json:"gender" binding:"required,oneof=male female other" example:"male"`
	IsEmailNotify *bool   `json:"is_email_notify" example:"true"`
}
