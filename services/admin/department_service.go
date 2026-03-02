package admin

import (
	"fmt"
	"phikhanh/models"
	adminRepo "phikhanh/repositories/admin"
	"phikhanh/utils"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// deptCodePattern - Chỉ cho phép đúng 3 chữ số sau "DP-"
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

	if rawCode == "" || name == "" {
		return &models.Department{
			Code:       rawCode,
			Name:       name,
			Address:    ctx.PostForm("address"),
			LeaderName: ctx.PostForm("leader_name"),
		}, utils.NewBadRequestError("Code and Name are required")
	}

	// Auto prefix DP- nếu user chỉ nhập 3 số
	code := normalizeCode(rawCode)

	// Validate format DP-XXX
	if !deptCodePattern.MatchString(code) {
		return &models.Department{
			Code:       rawCode,
			Name:       name,
			Address:    ctx.PostForm("address"),
			LeaderName: ctx.PostForm("leader_name"),
		}, utils.NewBadRequestError("Code must be in format DP-XXX (e.g. DP-001)")
	}

	department := &models.Department{
		Code:       code,
		Name:       name,
		Address:    ctx.PostForm("address"),
		LeaderName: ctx.PostForm("leader_name"),
	}

	return department, nil
}

// normalizeCode - Tự động thêm prefix DP-
func normalizeCode(raw string) string {
	if matched, _ := regexp.MatchString(`^\d{3}$`, raw); matched {
		return fmt.Sprintf("DP-%s", raw)
	}
	return raw
}
