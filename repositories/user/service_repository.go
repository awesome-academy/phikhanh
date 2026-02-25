package user

import (
	"phikhanh/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Lấy danh sách services với pagination và filter
func (r *ServiceRepository) GetServiceList(page, limit int, keyword, sector string, departmentID *uuid.UUID) ([]models.Service, int64, error) {
	var services []models.Service
	var total int64

	query := r.db.Model(&models.Service{}).Preload("Department")

	// Filter by keyword (name OR code)
	if keyword != "" {
		query = query.Where("(name ILIKE ? OR code ILIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Filter by sector
	if sector != "" {
		query = query.Where("sector = ?", sector)
	}

	// Filter by department_id
	if departmentID != nil {
		query = query.Where("department_id = ?", departmentID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&services).Error; err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

// Lấy chi tiết service theo ID
func (r *ServiceRepository) GetServiceByID(id uuid.UUID) (*models.Service, error) {
	var service models.Service
	err := r.db.Preload("Department").Where("id = ?", id).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}
