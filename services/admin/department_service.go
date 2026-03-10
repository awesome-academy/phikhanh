package admin

import (
	"strings"
	"time"

	adminDto "phikhanh/dto/admin"
	"phikhanh/models"
	adminRepo "phikhanh/repositories/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentService struct {
	repo *adminRepo.DepartmentRepository
}

func NewDepartmentService(repo *adminRepo.DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) GetAll() ([]models.Department, error) {
	departments, err := s.repo.FindAll()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}
	return departments, nil
}

func (s *DepartmentService) GetByID(id uuid.UUID) (*models.Department, error) {
	department, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFoundError("Department not found")
		}
		return nil, utils.NewInternalServerError(err)
	}
	return department, nil
}

// GetDetail - Lấy chi tiết department với formatted timestamps
func (s *DepartmentService) GetDetail(id uuid.UUID) (*adminDto.DepartmentDetail, error) {
	dept, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	detail := &adminDto.DepartmentDetail{
		ID:        dept.ID.String(),
		Code:      dept.Code,
		Name:      dept.Name,
		Address:   dept.Address,
		CreatedAt: dept.CreatedAt.Format(time.DateTime),
		UpdatedAt: dept.UpdatedAt.Format(time.DateTime),
	}

	if dept.Leader != nil {
		detail.LeaderName = dept.Leader.Name
	} else {
		detail.LeaderName = dept.LeaderName
	}

	return detail, nil
}

// GetAvailableManagers - Lấy danh sách managers để chọn làm leader
func (s *DepartmentService) GetAvailableManagers() ([]adminDto.ManagerOption, error) {
	managers, err := s.repo.GetAvailableManagers()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	result := make([]adminDto.ManagerOption, 0, len(managers))
	for _, m := range managers {
		result = append(result, adminDto.ManagerOption{
			ID:   m.ID.String(),
			Name: m.Name,
		})
	}
	return result, nil
}

func (s *DepartmentService) Create(department *models.Department) error {
	exists, err := s.repo.IsCodeExists(department.Code)
	if err != nil {
		return utils.NewInternalServerError(err)
	}
	if exists {
		return utils.NewBadRequestError("Department code already exists")
	}

	// Denormalize leader name
	if department.LeaderID != nil {
		managers, _ := s.repo.GetAvailableManagers()
		for _, m := range managers {
			if m.ID == *department.LeaderID {
				department.LeaderName = m.Name
				break
			}
		}
	}

	if err := s.repo.CreateWithLeader(department); err != nil {
		return utils.NewInternalServerError(err)
	}
	return nil
}

func (s *DepartmentService) Update(department *models.Department) error {
	current, err := s.repo.FindByID(department.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewNotFoundError("Department not found")
		}
		return utils.NewInternalServerError(err)
	}

	// Denormalize leader name
	department.LeaderName = ""
	if department.LeaderID != nil {
		managers, _ := s.repo.GetAvailableManagers()
		for _, m := range managers {
			if m.ID == *department.LeaderID {
				department.LeaderName = m.Name
				break
			}
		}
	}

	if err := s.repo.UpdateWithLeader(department, current.LeaderID); err != nil {
		return utils.NewInternalServerError(err)
	}
	return nil
}

func (s *DepartmentService) Delete(id uuid.UUID) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return utils.NewInternalServerError(err)
	}
	return nil
}

// BindForm - Parse form data bao gồm leader_id
func (s *DepartmentService) BindForm(ctx *gin.Context) (*models.Department, error) {
	code := strings.TrimSpace(ctx.PostForm("code"))
	name := strings.TrimSpace(ctx.PostForm("name"))
	address := strings.TrimSpace(ctx.PostForm("address"))
	leaderIDStr := strings.TrimSpace(ctx.PostForm("leader_id"))

	if code == "" {
		return &models.Department{}, utils.NewBadRequestError("Code is required")
	}
	if name == "" {
		return &models.Department{}, utils.NewBadRequestError("Name is required")
	}

	// Auto-prefix DP-
	if !strings.HasPrefix(code, "DP-") {
		code = "DP-" + code
	}

	dept := &models.Department{
		Code:    code,
		Name:    name,
		Address: address,
	}

	if leaderIDStr != "" {
		leaderID, err := uuid.Parse(leaderIDStr)
		if err != nil {
			return dept, utils.NewBadRequestError("Invalid leader ID")
		}
		dept.LeaderID = &leaderID
	}

	return dept, nil
}
