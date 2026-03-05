package admin

// ApplicationListItem - DTO item cho danh sách applications
type ApplicationListItem struct {
	ID            string
	Code          string
	ApplicantName string
	ServiceName   string
	Status        string
	SubmittedAt   string
}

// ApplicationListResult - Kết quả phân trang
type ApplicationListResult struct {
	Items       []ApplicationListItem
	CurrentPage int
	TotalPages  int
	TotalItems  int64
	Status      string
	HasPrev     bool
	HasNext     bool
}

// ApplicationAttachment - DTO attachment
type ApplicationAttachment struct {
	FileName string
	FilePath string
}

// ApplicationHistory - DTO history item
type ApplicationHistory struct {
	Date        string
	Action      string
	ActorName   string
	Note        string
	Description string
}

// ApplicationDetail - DTO chi tiết application
type ApplicationDetail struct {
	ID             string
	Code           string
	ApplicantName  string
	Email          string
	CitizenID      string
	Phone          string
	ServiceName    string
	ProcessingDays int
	Fee            *int
	Status         string
	SubmittedAt    string
	Histories      []ApplicationHistory
	Attachments    []ApplicationAttachment
}

// StaffMember - DTO staff member cho dropdown assign
type StaffMember struct {
	ID   string
	Name string
	Role string
}

// ActivityLogItem - DTO item cho danh sách hoạt động
type ActivityLogItem struct {
	ID        string
	AdminName string
	AdminRole string
	Action    string
	Module    string
	Details   string
	Status    string
	Timestamp string
}

// ActivityLogListResult - Kết quả phân trang
type ActivityLogListResult struct {
	Items       []ActivityLogItem
	CurrentPage int
	TotalPages  int
	TotalItems  int64
	HasPrev     bool
	HasNext     bool
}
