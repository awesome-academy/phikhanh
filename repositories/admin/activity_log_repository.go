package admin

import (
	"time"

	"phikhanh/models"

	"gorm.io/gorm"
)

type ActivityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

func (r *ActivityLogRepository) Create(log *models.SystemLog) error {
	return r.db.Create(log).Error
}

func (r *ActivityLogRepository) FindAllWithFilter(action string, offset, limit int) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64

	query := r.db.Model(&models.SystemLog{}).Preload("Actor")

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// DeleteOlderThan - Xóa logs cũ hơn n ngày
func (r *ActivityLogRepository) DeleteOlderThan(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := r.db.Where("created_at < ?", cutoff).Delete(&models.SystemLog{})
	return result.RowsAffected, result.Error
}

func (r *ActivityLogRepository) GetAvailableActions() []string {
	return []string{
		models.ActionLogin,
		models.ActionLogout,
		models.ActionSubmitApp,
		models.ActionUpdateApp,
		models.ActionDeleteApp,
		models.ActionUpdateStatus,
		models.ActionAssignStaff,
		models.ActionCreateService,
		models.ActionUpdateService,
		models.ActionDeleteService,
		models.ActionCreateUser,
		models.ActionUpdateUser,
		models.ActionDeleteUser,
		models.ActionCreateDept,
		models.ActionUpdateDept,
		models.ActionDeleteDept,
	}
}
