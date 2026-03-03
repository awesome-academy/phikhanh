package admin

import (
	adminDto "phikhanh/dto/admin"
	"phikhanh/models"
	adminRepo "phikhanh/repositories/admin"
	"phikhanh/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AvailableSectors - Danh sách sectors hợp lệ
var AvailableSectors = []string{
	"Health", "Land", "Construction", "Education",
	"Finance", "Transportation", "Agriculture", "Other",
}

type ServiceAdminService struct {
	repo *adminRepo.ServiceRepository
}

func NewServiceAdminService(repo *adminRepo.ServiceRepository) *ServiceAdminService {
	return &ServiceAdminService{repo: repo}
}

func (s *ServiceAdminService) GetAll() ([]models.Service, error) {
	services, err := s.repo.FindAll()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}
	return services, nil
}

func (s *ServiceAdminService) GetByID(id uuid.UUID) (*models.Service, error) {
	service, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFoundError("Service not found")
		}
		return nil, utils.NewInternalServerError(err)
	}
	return service, nil
}

// GetDetail - Lấy chi tiết service với formatted timestamps
func (s *ServiceAdminService) GetDetail(id uuid.UUID) (*adminDto.ServiceDetail, error) {
	svc, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	detail := &adminDto.ServiceDetail{
		ID:             svc.ID.String(),
		Code:           svc.Code,
		Name:           svc.Name,
		Sector:         svc.Sector,
		Description:    svc.Description,
		ProcessingDays: svc.ProcessingDays,
		Fee:            svc.Fee,
		CreatedAt:      svc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      svc.UpdatedAt.Format(time.RFC3339),
	}

	if svc.Department != nil {
		detail.DepartmentName = svc.Department.Name
	}

	return detail, nil
}

func (s *ServiceAdminService) Create(service *models.Service) error {
	if err := s.repo.Create(service); err != nil {
		return utils.ParseDBError(err)
	}
	return nil
}

func (s *ServiceAdminService) Update(service *models.Service) error {
	if err := s.repo.Update(service); err != nil {
		return utils.ParseDBError(err)
	}
	return nil
}

func (s *ServiceAdminService) Delete(id uuid.UUID) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return utils.ParseDBError(err)
	}
	return nil
}

func (s *ServiceAdminService) GetDepartments() ([]models.Department, error) {
	depts, err := s.repo.FindAllDepartments()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}
	return depts, nil
}

// BindForm - Parse form values thành Service model
func (s *ServiceAdminService) BindForm(ctx *gin.Context) (*models.Service, error) {
	service := &models.Service{
		Code:        ctx.PostForm("code"),
		Name:        ctx.PostForm("name"),
		Sector:      ctx.PostForm("sector"),
		Description: ctx.PostForm("description"),
	}

	if service.Code == "" || service.Name == "" {
		return service, utils.NewBadRequestError("Code and Name are required")
	}

	if days := ctx.PostForm("processing_days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil && d >= 0 {
			service.ProcessingDays = d
		}
	}

	if feeStr := ctx.PostForm("fee"); feeStr != "" {
		if f, err := strconv.Atoi(feeStr); err == nil && f >= 0 {
			service.Fee = &f
		}
	}

	if deptID := ctx.PostForm("department_id"); deptID != "" {
		if id, err := uuid.Parse(deptID); err == nil {
			service.DepartmentID = id
		}
	}

	return service, nil
}
