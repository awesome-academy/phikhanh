package models

import "github.com/google/uuid"

type Service struct {
	BaseModel
	Code           string    `gorm:"not null" json:"code"`
	Name           string    `gorm:"not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	Sector         string    `json:"sector"`
	DepartmentID   uuid.UUID `gorm:"type:uuid;not null" json:"department_id"`
	ProcessingDays int       `json:"processing_days"`
	Fee            *int      `json:"fee"`

	// Relations
	Department   *Department   `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Applications []Application `gorm:"foreignKey:ServiceID" json:"-"`
}
