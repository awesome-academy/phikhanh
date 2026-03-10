package admin

import (
	"phikhanh/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAllWithFilter - Lấy danh sách users (không include soft deleted) với filter role và pagination
func (r *UserRepository) FindAllWithFilter(role string, offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	if role != "" && role != "all" {
		query = query.Where("role = ?", role)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindByID - Lấy user theo ID
func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User

	if err := r.db.Preload("Department").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Create - Tạo user mới
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update - Cập nhật user dùng map để tránh GORM bỏ qua zero values
func (r *UserRepository) Update(user *models.User) error {
	updates := map[string]interface{}{
		"citizen_id":      user.CitizenID,
		"name":            user.Name,
		"email":           user.Email,
		"phone":           user.Phone,
		"address":         user.Address,
		"date_of_birth":   user.DateOfBirth,
		"gender":          user.Gender,
		"role":            user.Role,
		"department_id":   user.DepartmentID,
		"is_email_notify": user.IsEmailNotify,
	}

	return r.db.Model(user).Updates(updates).Error
}

// SoftDelete - Soft delete user (set deleted_at)
func (r *UserRepository) SoftDelete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

// GetUsersByRole - Lấy users theo role (không include soft deleted)
func (r *UserRepository) GetUsersByRole(role string) ([]models.User, error) {
	var users []models.User

	if err := r.db.Where("role = ?", role).Order("name ASC").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetDepartments - Lấy danh sách departments
func (r *UserRepository) GetDepartments() ([]models.Department, error) {
	var departments []models.Department

	if err := r.db.Order("name ASC").Find(&departments).Error; err != nil {
		return nil, err
	}

	return departments, nil
}

// IsCitizenIDExists - Kiểm tra citizen ID đã tồn tại
func (r *UserRepository) IsCitizenIDExists(citizenID string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("citizen_id = ?", citizenID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsEmailExists - Kiểm tra email đã tồn tại
func (r *UserRepository) IsEmailExists(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsCitizenIDExistsExcept - Kiểm tra citizen ID đã tồn tại (exclude current user)
func (r *UserRepository) IsCitizenIDExistsExcept(citizenID, userID string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("citizen_id = ? AND id != ?", citizenID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsEmailExistsExcept - Kiểm tra email đã tồn tại (exclude current user)
func (r *UserRepository) IsEmailExistsExcept(email, userID string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ? AND id != ?", email, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
