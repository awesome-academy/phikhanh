package admin

import (
	"fmt"
	"strconv"

	adminRepo "phikhanh/repositories/admin"
	"phikhanh/utils"
)

type ExportService struct {
	repo *adminRepo.ExportRepository
}

func NewExportService(repo *adminRepo.ExportRepository) *ExportService {
	return &ExportService{repo: repo}
}

// ExportCitizens - CitizenID, Name, Email, Phone, Total Applications
func (s *ExportService) ExportCitizens() ([][]string, error) {
	users, err := s.repo.GetCitizensWithAppCount()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	rows := make([][]string, 0, len(users)+1)
	rows = append(rows, []string{"Citizen ID", "Name", "Email", "Phone", "Total Applications"})

	for _, u := range users {
		rows = append(rows, []string{
			u.CitizenID,
			u.Name,
			u.Email,
			u.Phone,
			strconv.Itoa(len(u.Applications)),
		})
	}
	return rows, nil
}

// ExportApplications - Code, Citizen Name, Service Name, Status, Created At, Processing Days
func (s *ExportService) ExportApplications() ([][]string, error) {
	apps, err := s.repo.GetApplicationsWithDetails()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	rows := make([][]string, 0, len(apps)+1)
	rows = append(rows, []string{"Application Code", "Citizen Name", "Service Name", "Status", "Created At", "Processing Days"})

	for _, app := range apps {
		citizenName := ""
		if app.User != nil {
			citizenName = app.User.Name
		}

		serviceName := ""
		processingDays := ""
		if app.Service != nil {
			serviceName = app.Service.Name
			processingDays = strconv.Itoa(app.Service.ProcessingDays)
		}

		rows = append(rows, []string{
			app.Code,
			citizenName,
			serviceName,
			string(app.Status),
			app.CreatedAt.Format("2006-01-02 15:04:05"),
			processingDays,
		})
	}
	return rows, nil
}

// ExportServices - Code, Name, Sector, Department Name, Fee, Processing Days
func (s *ExportService) ExportServices() ([][]string, error) {
	services, err := s.repo.GetServicesWithDepartment()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	rows := make([][]string, 0, len(services)+1)
	rows = append(rows, []string{"Service Code", "Service Name", "Sector", "Department Name", "Fee (VND)", "Processing Days"})

	for _, svc := range services {
		deptName := ""
		if svc.Department != nil {
			deptName = svc.Department.Name
		}

		fee := "Free"
		if svc.Fee != nil && *svc.Fee > 0 {
			fee = fmt.Sprintf("%d", *svc.Fee)
		}

		rows = append(rows, []string{
			svc.Code,
			svc.Name,
			svc.Sector,
			deptName,
			fee,
			strconv.Itoa(svc.ProcessingDays),
		})
	}
	return rows, nil
}

// ExportDepartments - Code, Name, Address, Leader Name
func (s *ExportService) ExportDepartments() ([][]string, error) {
	departments, err := s.repo.GetDepartments()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	rows := make([][]string, 0, len(departments)+1)
	rows = append(rows, []string{"Department Code", "Name", "Address", "Leader Name"})

	for _, dept := range departments {
		leaderName := dept.LeaderName
		if dept.Leader != nil {
			leaderName = dept.Leader.Name
		}

		rows = append(rows, []string{
			dept.Code,
			dept.Name,
			dept.Address,
			leaderName,
		})
	}
	return rows, nil
}

// ExportStaff - Staff ID, Name, Email, Department Name, Role
func (s *ExportService) ExportStaff() ([][]string, error) {
	staff, err := s.repo.GetStaffWithDepartment()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	rows := make([][]string, 0, len(staff)+1)
	rows = append(rows, []string{"Staff ID", "Name", "Email", "Department Name", "Role"})

	for _, u := range staff {
		deptName := "-"
		if u.Department != nil {
			deptName = u.Department.Name
		}

		rows = append(rows, []string{
			u.CitizenID,
			u.Name,
			u.Email,
			deptName,
			string(u.Role),
		})
	}
	return rows, nil
}
