package user

import (
	"phikhanh/models"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// Tạo user mới
func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// Tìm user theo citizen_id
func (r *AuthRepository) FindByCitizenID(citizenID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("citizen_id = ?", citizenID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Tìm user theo email
func (r *AuthRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Kiểm tra citizen_id đã tồn tại
func (r *AuthRepository) IsCitizenIDExists(citizenID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("citizen_id = ?", citizenID).Count(&count).Error
	return count > 0, err
}

// Kiểm tra email đã tồn tại
func (r *AuthRepository) IsEmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
