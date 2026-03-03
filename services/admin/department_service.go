package admin

import (
	"fmt"
	adminDto "phikhanh/dto/admin"
	"phikhanh/models"
	adminRepo "phikhanh/repositories/admin"
	"phikhanh/utils"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var deptCodePattern = regexp.MustCompile(`^DP-\d{3}$`)

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

	return &adminDto.DepartmentDetail{
		ID:         dept.ID.String(),
		Code:       dept.Code,
		Name:       dept.Name,
		Address:    dept.Address,
		LeaderName: dept.LeaderName,
		CreatedAt:  dept.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  dept.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *DepartmentService) Create(department *models.Department) error {
	if err := s.repo.Create(department); err != nil {
		return utils.ParseDBError(err)
	}
	return nil
}

func (s *DepartmentService) Update(department *models.Department) error {
	if err := s.repo.Update(department); err != nil {
		return utils.ParseDBError(err)
	}
	return nil
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

// BindForm - Parse và validate form values thành Department model
func (s *DepartmentService) BindForm(ctx *gin.Context) (*models.Department, error) {
	rawCode := ctx.PostForm("code")
	name := ctx.PostForm("name")

	dept := &models.Department{
		Code:       rawCode,
		Name:       name,
		Address:    ctx.PostForm("address"),
		LeaderName: ctx.PostForm("leader_name"),
	}

	if rawCode == "" || name == "" {
		return dept, utils.NewBadRequestError("Code and Name are required")
	}

	code := normalizeDeptCode(rawCode)
	if !deptCodePattern.MatchString(code) {
		return dept, utils.NewBadRequestError("Code must be in format DP-XXX (e.g. DP-001)")
	}

	dept.Code = code
	return dept, nil
}

// normalizeDeptCode - Tự động thêm prefix DP- nếu chưa có
func normalizeDeptCode(raw string) string {
	if matched, _ := regexp.MatchString(`^\d{3}$`, raw); matched {
		return fmt.Sprintf("DP-%s", raw)
	}
	return raw
}
