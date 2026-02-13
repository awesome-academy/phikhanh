package user

// Request đăng ký tài khoản
type RegisterRequest struct {
	CitizenID   string `json:"citizen_id" binding:"required,citizen_id" example:"001234567890"`
	Password    string `json:"password" binding:"required,strong_password" example:"Password@123"`
	Name        string `json:"name" binding:"required" example:"Nguyễn Văn A"`
	Email       string `json:"email" binding:"required,email" example:"nguyenvana@example.com"`
	Phone       string `json:"phone" binding:"required,vn_phone" example:"0901234567"`
	Address     string `json:"address" example:"123 Đường ABC, Quận 1, TP.HCM"`
	DateOfBirth string `json:"date_of_birth" binding:"omitempty,past_date" example:"1990-01-15"`
	Gender      string `json:"gender" binding:"required,oneof=male female other" example:"male"`
}

// Request đăng nhập
type LoginRequest struct {
	CitizenID string `json:"citizen_id" binding:"required" example:"001234567890"`
	Password  string `json:"password" binding:"required" example:"password123"`
}

// Response sau khi đăng nhập thành công
type LoginResponse struct {
	Token string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserInfo `json:"user"`
}

// Thông tin user cơ bản
type UserInfo struct {
	ID        string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CitizenID string `json:"citizen_id" example:"001234567890"`
	Name      string `json:"name" example:"Nguyễn Văn A"`
	Email     string `json:"email" example:"nguyenvana@example.com"`
	Phone     string `json:"phone" example:"0901234567"`
	Role      string `json:"role" example:"citizen"`
}
