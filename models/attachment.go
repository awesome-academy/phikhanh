package models

import "github.com/google/uuid"

type AttachmentType string

const (
	AttachmentTypeOriginal   AttachmentType = "original"
	AttachmentTypeSupplement AttachmentType = "supplement"
)

type Attachment struct {
	BaseModel
	ApplicationID uuid.UUID `gorm:"type:uuid;not null;index" json:"application_id"`
	FileName      string    `gorm:"not null" json:"file_name"`
	FilePath      string    `gorm:"not null" json:"file_path"`
	Type          string    `gorm:"type:varchar(50);default:'original'" json:"type"` // "original" | "supplement"

	// Relations
	Application *Application `gorm:"foreignKey:ApplicationID" json:"-"`
}
