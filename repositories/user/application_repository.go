package user

import (
	userDto "phikhanh/dto/user"
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

// FindMyApplications - Lấy danh sách hồ sơ của user với pagination và filter
func (r *ApplicationRepository) FindMyApplications(
	userID uuid.UUID,
	req userDto.MyAppListRequest,
) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	query := r.db.Model(&models.Application{}).
		Where("user_id = ?", userID)

	// Filter by status nếu có
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Count total trước khi pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Chỉ preload Service để lấy Name, không preload Attachments/Histories
	offset := (req.Page - 1) * req.Limit
	if err := query.
		Preload("Service", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name") // Chỉ lấy các field cần thiết
		}).
		Order("created_at DESC").
		Offset(offset).
		Limit(req.Limit).
		Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}
