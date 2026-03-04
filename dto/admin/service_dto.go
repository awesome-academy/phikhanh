package admin

// ServiceDetail - DTO chi tiết service cho SSR
type ServiceDetail struct {
	ID             string
	Code           string
	Name           string
	Sector         string
	Description    string
	ProcessingDays int
	Fee            *int
	DepartmentName string
	CreatedAt      string
	UpdatedAt      string
}
