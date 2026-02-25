package user

import (
	"phikhanh/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Lấy thông tin user theo ID
func (r *ProfileRepository) FindByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Cập nhật thông tin user
func (r *ProfileRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}
