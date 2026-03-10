package models

import (
	"time"

	"github.com/google/uuid"
)

// Action constants
const (
	ActionLogin         = "LOGIN"
	ActionLogout        = "LOGOUT"
	ActionSubmitApp     = "SUBMIT_APP"
	ActionUpdateApp     = "UPDATE_APP"
	ActionDeleteApp     = "DELETE_APP"
	ActionUpdateStatus  = "UPDATE_STATUS"
	ActionAssignStaff   = "ASSIGN_STAFF"
	ActionCreateService = "CREATE_SERVICE"
	ActionUpdateService = "UPDATE_SERVICE"
	ActionDeleteService = "DELETE_SERVICE"
	ActionCreateUser    = "CREATE_USER"
	ActionUpdateUser    = "UPDATE_USER"
	ActionDeleteUser    = "DELETE_USER"
	ActionCreateDept    = "CREATE_DEPT"
	ActionUpdateDept    = "UPDATE_DEPT"
	ActionDeleteDept    = "DELETE_DEPT"
)

// Model đại diện cho bảng system_logs
type SystemLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ActorID     *uuid.UUID `gorm:"type:uuid;index"`
	Action      string     `gorm:"type:varchar(100);not null;index"`
	TargetID    string     `gorm:"type:varchar(255)"`
	Description string     `gorm:"type:text"`
	IPAddress   string     `gorm:"type:varchar(45)"`
	CreatedAt   time.Time  `gorm:"index"`

	// Relations
	Actor *User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

func (SystemLog) TableName() string { return "system_logs" }
