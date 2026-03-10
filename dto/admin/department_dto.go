package admin

// DepartmentDetail - DTO chi tiết department cho SSR
type DepartmentDetail struct {
	ID         string
	Code       string
	Name       string
	Address    string
	LeaderName string
	CreatedAt  string
	UpdatedAt  string
}

// ManagerOption - DTO chi tiết manager cho SSR
type ManagerOption struct {
	ID   string
	Name string
}
