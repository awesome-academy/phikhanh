package models

import "github.com/google/uuid"

// Model đại diện cho bảng notifications
type Notification struct {
	BaseModel
	UserID  uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title   string    `gorm:"not null" json:"title"`
	Content string    `gorm:"type:text" json:"content"`
	IsRead  bool      `gorm:"default:false" json:"is_read"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"-"`
}
