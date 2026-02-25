package user

import (
	"phikhanh/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// IsServiceExists - Kiểm tra service có tồn tại không
func (r *ApplicationRepository) IsServiceExists(serviceID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Service{}).Where("id = ?", serviceID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateWithTransaction - Tạo application, attachments và history trong 1 transaction
func (r *ApplicationRepository) CreateWithTransaction(
	app *models.Application,
	attachments []models.Attachment,
	history *models.ApplicationHistory,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Insert application
		if err := tx.Create(app).Error; err != nil {
			return err
		}

		// 2. Insert attachments nếu có
		if len(attachments) > 0 {
			for i := range attachments {
				attachments[i].ApplicationID = app.ID
			}
			if err := tx.Create(&attachments).Error; err != nil {
				return err
			}
		}

		// 3. Insert history
		history.ApplicationID = app.ID
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		return nil
	})
}
