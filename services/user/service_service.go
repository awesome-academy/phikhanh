package user

import (
	"math"

	userDto "phikhanh/dto/user"
	userRepo "phikhanh/repositories/user"
	"phikhanh/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceService struct {
	repo *userRepo.ServiceRepository
}

func NewServiceService(repo *userRepo.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

// Lấy danh sách services
func (s *ServiceService) GetServiceList(req userDto.ServiceListRequest) (*userDto.ServiceListResponse, error) {
	// Set default values nếu không được truyền
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Parse department_id - trả về lỗi nếu không hợp lệ
	var departmentID *uuid.UUID
	if req.DepartmentID != "" {
		id, err := uuid.Parse(req.DepartmentID)
		if err != nil {
			return nil, utils.NewBadRequestError("Invalid department_id format")
		}
		departmentID = &id
	}

	// Get services from repository
	services, total, err := s.repo.GetServiceList(req.Page, req.Limit, req.Keyword, req.Sector, departmentID)
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	// Map to DTO
	items := make([]userDto.ServiceListItem, 0, len(services))
	for _, service := range services {
		item := userDto.ServiceListItem{
			ID:             service.ID.String(),
			Name:           service.Name,
			Code:           service.Code,
			Sector:         service.Sector,
			Fee:            service.Fee,
			ProcessingDays: service.ProcessingDays,
		}
		if service.Department != nil {
			item.DepartmentName = service.Department.Name
		}
		items = append(items, item)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &userDto.ServiceListResponse{
		Items:      items,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

// Lấy chi tiết service
func (s *ServiceService) GetServiceDetail(id uuid.UUID) (*userDto.ServiceDetailResponse, error) {
	service, err := s.repo.GetServiceByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFoundError("Service not found")
		}
		return nil, utils.NewInternalServerError(err)
	}

	response := &userDto.ServiceDetailResponse{
		ID:             service.ID.String(),
		Name:           service.Name,
		Code:           service.Code,
		Description:    service.Description,
		Sector:         service.Sector,
		Fee:            service.Fee,
		ProcessingDays: service.ProcessingDays,
		CreatedAt:      service.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if service.Department != nil {
		response.Department = userDto.DepartmentInfo{
			ID:      service.Department.ID.String(),
			Name:    service.Department.Name,
			Code:    service.Department.Code,
			Address: service.Department.Address,
		}
	}

	return response, nil
}
