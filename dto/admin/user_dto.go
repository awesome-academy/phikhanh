package admin

type UserListItem struct {
	ID        string
	CitizenID string
	Name      string
	Email     string
	Role      string
	Status    string
}

type UserListResult struct {
	Items       []UserListItem
	CurrentPage int
	TotalPages  int
	TotalItems  int64
	Role        string
	HasPrev     bool
	HasNext     bool
}

type UserDetail struct {
	ID             string
	CitizenID      string
	Name           string
	Email          string
	Phone          string
	Address        string
	DateOfBirth    string
	Gender         string
	Role           string
	DepartmentID   string
	DepartmentName string
	CreatedAt      string
	UpdatedAt      string
}

type CreateUserRequest struct {
	CitizenID    string `form:"citizen_id" binding:"required,citizen_id"`
	Name         string `form:"name" binding:"required"`
	Email        string `form:"email" binding:"required,email"`
	Password     string `form:"password" binding:"required,strong_password"`
	Role         string `form:"role" binding:"required,oneof=citizen staff manager admin"`
	Phone        string `form:"phone" binding:"omitempty,vn_phone"`
	Address      string `form:"address"`
	DateOfBirth  string `form:"date_of_birth" binding:"omitempty,past_date"`
	Gender       string `form:"gender" binding:"omitempty,oneof=male female other"`
	DepartmentID string `form:"department_id"`
}

type UpdateUserRequest struct {
	CitizenID    string `form:"citizen_id" binding:"required,citizen_id"`
	Name         string `form:"name" binding:"required"`
	Email        string `form:"email" binding:"required,email"`
	Role         string `form:"role" binding:"required,oneof=citizen staff manager admin"`
	Phone        string `form:"phone" binding:"omitempty,vn_phone"`
	Address      string `form:"address"`
	DateOfBirth  string `form:"date_of_birth" binding:"omitempty,past_date"`
	Gender       string `form:"gender" binding:"omitempty,oneof=male female other"`
	DepartmentID string `form:"department_id"`
}

type RoleOption struct {
	Value string
	Label string
}

type GenderOption struct {
	Value string
	Label string
}

type DepartmentOption struct {
	ID   string
	Name string
}

// UserFormData - DTO để render form (tránh dùng models.User trực tiếp trong template)
type UserFormData struct {
	CitizenID    string
	Name         string
	Email        string
	Phone        string
	Address      string
	DateOfBirth  string // "2006-01-02" format
	Gender       string
	Role         string
	DepartmentID string
}
