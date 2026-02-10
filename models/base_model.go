package models

import (
	"time"

	"gorm.io/gorm"
)

// Struct cơ bản chứa các trường chung cho tất cả models
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
