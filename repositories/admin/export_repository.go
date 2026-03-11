package admin

import (
	"phikhanh/models"

	"gorm.io/gorm"
)

type ExportRepository struct {
	db *gorm.DB
}

func NewExportRepository(db *gorm.DB) *ExportRepository {
	return &ExportRepository{db: db}
}

func (r *ExportRepository) GetCitizensWithAppCount() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ?", models.RoleCitizen).
		Preload("Applications").
		Order("created_at DESC").
		Find(&users).Error
	return users, err
}

func (r *ExportRepository) GetApplicationsWithDetails() ([]models.Application, error) {
	var apps []models.Application
	err := r.db.Preload("User").
		Preload("Service").
		Order("created_at DESC").
		Find(&apps).Error
	return apps, err
}

func (r *ExportRepository) GetServicesWithDepartment() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Preload("Department").
		Order("created_at DESC").
		Find(&services).Error
	return services, err
}

func (r *ExportRepository) GetDepartments() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Preload("Leader").
		Order("code ASC").
		Find(&departments).Error
	return departments, err
}

func (r *ExportRepository) GetStaffWithDepartment() ([]models.User, error) {
	var staff []models.User
	err := r.db.Where("role = ?", models.RoleStaff).
		Preload("Department").
		Order("name ASC").
		Find(&staff).Error
	return staff, err
}
