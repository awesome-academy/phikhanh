package admin

import (
	"phikhanh/models"

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

// FindByEmail - Tìm user theo email (chỉ Staff, Manager, Admin)
func (r *AdminRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where(
		"email = ? AND role IN (?, ?, ?)",
		email,
		models.RoleStaff,
		models.RoleManager,
		models.RoleAdmin,
	).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
