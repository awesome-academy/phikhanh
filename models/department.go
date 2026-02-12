package models

// Model đại diện cho bảng departments
type Department struct {
	BaseModel
	Code       string `gorm:"unique;not null" json:"code"`
	Name       string `gorm:"not null" json:"name"`
	Address    string `json:"address"`
	LeaderName string `json:"leader_name"`

	// Relations
	Users    []User    `gorm:"foreignKey:DepartmentID" json:"-"`
	Services []Service `gorm:"foreignKey:DepartmentID" json:"-"`
}
