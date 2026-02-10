package admin

import (
	"gorm.io/gorm"
)

// AdminRepository - Repository xử lý truy vấn dữ liệu cho admin
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository - Khởi tạo AdminRepository mới
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}
