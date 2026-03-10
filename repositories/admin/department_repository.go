package admin

import (
	"phikhanh/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) FindAll() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Order("code ASC").Find(&departments).Error
	return departments, err
}

// GetAvailableManagers - Lấy danh sách users có role=manager để làm leader
func (r *DepartmentRepository) GetAvailableManagers() ([]models.User, error) {
	var managers []models.User
	if err := r.db.Where("role = ? AND deleted_at IS NULL", models.RoleManager).
		Order("name ASC").
		Find(&managers).Error; err != nil {
		return nil, err
	}
	return managers, nil
}

func (r *DepartmentRepository) FindByID(id uuid.UUID) (*models.Department, error) {
	var dept models.Department
	if err := r.db.Preload("Leader").First(&dept, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

func (r *DepartmentRepository) CreateWithLeader(dept *models.Department) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(dept).Error; err != nil {
			return err
		}
		// Assign department cho manager được chọn làm leader
		if dept.LeaderID != nil {
			if err := tx.Model(&models.User{}).
				Where("id = ? AND role = ?", dept.LeaderID, models.RoleManager).
				Update("department_id", dept.ID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *DepartmentRepository) Update(department *models.Department) error {
	return r.db.Save(department).Error
}

func (r *DepartmentRepository) UpdateWithLeader(dept *models.Department, oldLeaderID *uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(dept).Updates(map[string]interface{}{
			"code":        dept.Code,
			"name":        dept.Name,
			"address":     dept.Address,
			"leader_id":   dept.LeaderID,
			"leader_name": dept.LeaderName,
		}).Error; err != nil {
			return err
		}

		// Clear DepartmentID của leader cũ (nếu đổi leader)
		leaderChanged := (oldLeaderID == nil && dept.LeaderID != nil) ||
			(oldLeaderID != nil && dept.LeaderID == nil) ||
			(oldLeaderID != nil && dept.LeaderID != nil && *oldLeaderID != *dept.LeaderID)

		if leaderChanged && oldLeaderID != nil {
			if err := tx.Model(&models.User{}).
				Where("id = ?", oldLeaderID).
				Update("department_id", nil).Error; err != nil {
				return err
			}
		}

		// Assign department cho leader mới
		if dept.LeaderID != nil {
			if err := tx.Model(&models.User{}).
				Where("id = ? AND role = ?", dept.LeaderID, models.RoleManager).
				Update("department_id", dept.ID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *DepartmentRepository) IsCodeExists(code string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Department{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *DepartmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Department{}, "id = ?", id).Error
}
