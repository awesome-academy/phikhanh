package admin

import (
	"fmt"
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

// FindAllWithFilter - Lấy danh sách applications với filter và pagination
func (r *ApplicationRepository) FindAllWithFilter(status string, offset, limit int) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	query := r.db.Model(&models.Application{}).
		Preload("User").
		Preload("Service")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

// FindByIDWithDetails - Lấy chi tiết application với tất cả related data
func (r *ApplicationRepository) FindByIDWithDetails(id string) (*models.Application, error) {
	var application models.Application
	err := r.db.
		Preload("User").
		Preload("Service").
		Preload("Attachments").
		Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Histories.Actor").
		First(&application, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

// FindByIDWithDetailsAndDescription - Lấy chi tiết application với tất cả related data
func (r *ApplicationRepository) FindByIDWithDetailsAndDescription(id string) (*models.Application, error) {
	var application models.Application
	err := r.db.
		Preload("User").
		Preload("Service").
		Preload("Attachments").
		Preload("Histories", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Histories.Actor").
		Preload("AssignedStaff").
		First(&application, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

// GetStaffNameByID - Lấy tên staff theo ID
func (r *ApplicationRepository) GetStaffNameByID(staffID string) (string, error) {
	var user models.User
	err := r.db.Select("name").Where("id = ?", staffID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return user.Name, nil
}

// ProcessAndAssignWithHistory - Cập nhật status, assign staff, và thêm history trong transaction
func (r *ApplicationRepository) ProcessAndAssignWithHistory(appID string, newStatus string, assignedStaffID *string, note string, actorID string) error {
	// Parse appID từ string sang UUID
	appUUID, err := uuid.Parse(appID)
	if err != nil {
		return err
	}

	// Parse actorID từ string sang UUID
	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update application status
		updateData := map[string]interface{}{
			"status": newStatus,
		}

		// Chỉ update assigned_staff_id nếu được cung cấp
		if assignedStaffID != nil && *assignedStaffID != "" {
			staffUUID, err := uuid.Parse(*assignedStaffID)
			if err != nil {
				return err
			}
			updateData["assigned_staff_id"] = staffUUID
		}

		if err := tx.Model(&models.Application{}).Where("id = ?", appUUID).Updates(updateData).Error; err != nil {
			return err
		}

		// Insert history record - dùng UUID trực tiếp, không convert sang string
		history := &models.ApplicationHistory{
			ApplicationID: appUUID,
			ActorID:       actorUUID,
			Action:        newStatus,
			Note:          note,
		}

		if err := tx.Create(history).Error; err != nil {
			return err
		}

		return nil
	})
}

// ProcessAndAssignWithHistoryV2 - Version 2 với description chi tiết
func (r *ApplicationRepository) ProcessAndAssignWithHistoryV2(appID string, oldStatus string, newStatus string, assignedStaffID *string, assignedStaffName *string, note string, actorID string, description string) error {
	appUUID, err := uuid.Parse(appID)
	if err != nil {
		return err
	}

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		updateData := map[string]interface{}{
			"status": newStatus,
		}

		if assignedStaffID != nil && *assignedStaffID != "" {
			staffUUID, err := uuid.Parse(*assignedStaffID)
			if err != nil {
				return err
			}
			updateData["assigned_staff_id"] = staffUUID
		}

		if err := tx.Model(&models.Application{}).Where("id = ?", appUUID).Updates(updateData).Error; err != nil {
			return err
		}

		// Insert history record với description chi tiết
		history := &models.ApplicationHistory{
			ApplicationID: appUUID,
			ActorID:       actorUUID,
			Action:        fmt.Sprintf("%s → %s", oldStatus, newStatus),
			Note:          description,
		}

		if err := tx.Create(history).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetAvailableStaff - Lấy danh sách staff có thể assign
func (r *ApplicationRepository) GetAvailableStaff() ([]models.User, error) {
	var staff []models.User

	if err := r.db.Where("role = ?", "staff").Order("name ASC").Find(&staff).Error; err != nil {
		return nil, err
	}

	return staff, nil
}

// FindAllWithFilterAndAssignment - Lấy danh sách applications với filter status và assignment
func (r *ApplicationRepository) FindAllWithFilterAndAssignment(status string, assignedToUserID *string, offset, limit int) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	query := r.db.Model(&models.Application{}).
		Preload("User").
		Preload("Service")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Nếu assignedToUserID được cung cấp, filter applications assigned to user này
	if assignedToUserID != nil && *assignedToUserID != "" {
		query = query.Where("assigned_staff_id = ?", *assignedToUserID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}
