package admin

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

func (r *ServiceRepository) FindAll() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Preload("Department").Order("code ASC").Find(&services).Error
	return services, err
}

func (r *ServiceRepository) FindByID(id uuid.UUID) (*models.Service, error) {
	var service models.Service
	err := r.db.Preload("Department").First(&service, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) Create(service *models.Service) error {
	return r.db.Create(service).Error
}

func (r *ServiceRepository) Update(service *models.Service) error {
	return r.db.Save(service).Error
}

func (r *ServiceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Service{}, "id = ?", id).Error
}

func (r *ServiceRepository) FindAllDepartments() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Select("id, name").Order("name ASC").Find(&departments).Error
	return departments, err
}
