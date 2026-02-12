package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string
type Gender string

const (
	RoleCitizen UserRole = "citizen"
	RoleStaff   UserRole = "staff"
	RoleManager UserRole = "manager"
	RoleAdmin   UserRole = "admin"

	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// Model đại diện cho bảng users
type User struct {
	BaseModel
	CitizenID     string     `gorm:"uniqueIndex:idx_users_citizen_id;not null" json:"citizen_id"`
	PasswordHash  string     `gorm:"not null" json:"-"`
	Name          string     `gorm:"not null" json:"name"`
	Email         string     `gorm:"uniqueIndex:idx_users_email;not null" json:"email"`
	Phone         string     `json:"phone"`
	Address       string     `json:"address"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	Gender        Gender     `gorm:"type:varchar(10)" json:"gender"`
	Role          UserRole   `gorm:"type:varchar(20);default:'citizen'" json:"role"`
	DepartmentID  *uuid.UUID `gorm:"type:uuid" json:"department_id"`
	IsEmailNotify bool       `gorm:"default:true" json:"is_email_notify"`

	// Relations
	Department   *Department   `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Applications []Application `gorm:"foreignKey:UserID" json:"-"`
}
