package admin

import (
	"regexp"
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

	// Validate và denormalize leader
	if department.LeaderID != nil {
		leaderName, err := s.validateAndGetLeaderName(*department.LeaderID)
		if err != nil {
			return err
		}
		department.LeaderName = leaderName
	}

	if err := s.repo.CreateWithLeader(department); err != nil {
		return utils.ParseDBError(err)
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

	// Validate và denormalize leader
	department.LeaderName = ""
	if department.LeaderID != nil {
		leaderName, err := s.validateAndGetLeaderName(*department.LeaderID)
		if err != nil {
			return err
		}
		department.LeaderName = leaderName
	}

	if err := s.repo.UpdateWithLeader(department, current.LeaderID); err != nil {
		return utils.ParseDBError(err)
	}
	return nil
}

// validateAndGetLeaderName - Validate leader_id phải là manager thực sự tồn tại
// Trả về tên leader nếu hợp lệ, error nếu không phải manager hoặc không tồn tại
func (s *DepartmentService) validateAndGetLeaderName(leaderID uuid.UUID) (string, error) {
	managers, err := s.repo.GetAvailableManagers()
	if err != nil {
		return "", utils.NewInternalServerError(err)
	}

	for _, m := range managers {
		if m.ID == leaderID {
			return m.Name, nil
		}
	}

	// leaderID không tồn tại trong danh sách managers → tampering hoặc invalid
	return "", utils.NewBadRequestError("Selected leader is not a valid manager")
}

func (s *DepartmentService) Delete(id uuid.UUID) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return utils.ParseDBError(err)
	}
	return nil
}

// BindForm - Parse và validate form data bao gồm leader_id
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

	code = strings.ToUpper(code)
	rawDigits := strings.TrimPrefix(code, "DP-")

	if !regexp.MustCompile(`^\d{3}$`).MatchString(rawDigits) {
		return &models.Department{}, utils.NewBadRequestError("Code must be exactly 3 digits (e.g. 001). Final format: DP-001")
	}

	finalCode := "DP-" + rawDigits

	dept := &models.Department{
		Code:    finalCode,
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
