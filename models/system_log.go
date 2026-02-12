package models

import "github.com/google/uuid"

// Model đại diện cho bảng system_logs
type SystemLog struct {
	BaseModelWithoutSoftDelete
	ActorID     *uuid.UUID `gorm:"type:uuid;index" json:"actor_id"`
	Action      string     `gorm:"not null" json:"action"`
	TargetID    string     `json:"target_id"`
	Description string     `gorm:"type:text" json:"description"`
	IPAddress   string     `json:"ip_address"`

	// Relations
	Actor *User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}
