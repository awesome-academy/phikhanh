package models

import "github.com/google/uuid"

type ApplicationStatus string

const (
	StatusReceived           ApplicationStatus = "Received"
	StatusProcessing         ApplicationStatus = "Processing"
	StatusSupplementRequired ApplicationStatus = "Supplement_Required"
	StatusApproved           ApplicationStatus = "Approved"
	StatusRejected           ApplicationStatus = "Rejected"
)

// Model đại diện cho bảng applications
type Application struct {
	BaseModel
	Code            string            `gorm:"uniqueIndex:idx_applications_code;not null" json:"code"` // HS001
	UserID          uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	ServiceID       uuid.UUID         `gorm:"type:uuid;not null" json:"service_id"`
	AssignedStaffID *uuid.UUID        `gorm:"type:uuid" json:"assigned_staff_id"`
	Status          ApplicationStatus `gorm:"type:varchar(30);default:'Received'" json:"status"`

	// Relations
	User          *User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Service       *Service             `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	AssignedStaff *User                `gorm:"foreignKey:AssignedStaffID" json:"assigned_staff,omitempty"`
	Attachments   []Attachment         `gorm:"foreignKey:ApplicationID" json:"-"`
	Histories     []ApplicationHistory `gorm:"foreignKey:ApplicationID" json:"-"`
}
