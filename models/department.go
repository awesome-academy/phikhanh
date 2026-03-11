package models

import "github.com/google/uuid"

// Model đại diện cho bảng departments
type Department struct {
	BaseModel
	Code       string     `gorm:"not null" json:"code"`
	Name       string     `gorm:"not null" json:"name"`
	Address    string     `json:"address"`
	LeaderID   *uuid.UUID `gorm:"type:uuid" json:"leader_id"`
	LeaderName string     `json:"leader_name"` // denormalized for display

	// Relations
	Leader   *User     `gorm:"foreignKey:LeaderID;constraint:-" json:"leader,omitempty"`
	Users    []User    `gorm:"foreignKey:DepartmentID" json:"-"`
	Services []Service `gorm:"foreignKey:DepartmentID" json:"-"`
}
