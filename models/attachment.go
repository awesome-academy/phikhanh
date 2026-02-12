package models

import "github.com/google/uuid"

type AttachmentType string

const (
	AttachmentTypeRequest  AttachmentType = "request"
	AttachmentTypeResponse AttachmentType = "response"
)

// Model đại diện cho bảng attachments
type Attachment struct {
	BaseModel
	ApplicationID uuid.UUID      `gorm:"type:uuid;not null" json:"application_id"`
	FileName      string         `gorm:"not null" json:"file_name"`
	FilePath      string         `gorm:"not null" json:"file_path"`
	Type          AttachmentType `gorm:"type:varchar(20)" json:"type"`

	// Relations
	Application *Application `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
}
