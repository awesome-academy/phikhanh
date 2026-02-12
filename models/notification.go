package models

import "github.com/google/uuid"

type NotificationType string

const (
	NotificationTypeSystem      NotificationType = "system"
	NotificationTypeApplication NotificationType = "application"
)

// Model đại diện cho bảng notifications
type Notification struct {
	BaseModel
	UserID  uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	Title   string           `gorm:"not null" json:"title"`
	Message string           `gorm:"type:text" json:"message"`
	Type    NotificationType `gorm:"type:varchar(20)" json:"type"`
	IsRead  bool             `gorm:"default:false" json:"is_read"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
