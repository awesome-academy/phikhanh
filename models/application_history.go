package models

import "github.com/google/uuid"

// Model đại diện cho bảng application_histories
type ApplicationHistory struct {
	BaseModelWithoutSoftDelete
	ApplicationID uuid.UUID `gorm:"type:uuid;not null;index" json:"application_id"`
	ActorID       uuid.UUID `gorm:"type:uuid;not null" json:"actor_id"`
	Action        string    `gorm:"not null" json:"action"`
	Note          string    `gorm:"type:text" json:"note"`

	// Relations
	Application *Application `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	Actor       *User        `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}
