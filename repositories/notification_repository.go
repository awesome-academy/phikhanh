package repositories

import (
	"phikhanh/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *models.Notification) error {
	return r.db.Create(n).Error
}

// FindUserWithEmailPref - Fetch user email + IsEmailNotify (chỉ select fields cần thiết)
func (r *NotificationRepository) FindUserWithEmailPref(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Select("id", "email", "name", "is_email_notify").
		First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUserID - Lấy danh sách notifications của user, mới nhất trước
func (r *NotificationRepository) FindByUserID(userID uuid.UUID, offset, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// MarkAsRead - Đánh dấu notification đã đọc (chỉ cho phép nếu thuộc về user)
func (r *NotificationRepository) MarkAsRead(notificationID, userID uuid.UUID) error {
	result := r.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// MarkAllAsRead - Đánh dấu tất cả notifications của user đã đọc
func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

// CountUnread - Đếm số notifications chưa đọc
func (r *NotificationRepository) CountUnread(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}
