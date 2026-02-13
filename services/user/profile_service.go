package user

import (
	"time"

	"phikhanh/models"
	userRepo "phikhanh/repositories/user"
	"phikhanh/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileService struct {
	repo *userRepo.ProfileRepository
}

func NewProfileService(repo *userRepo.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

// Lấy thông tin profile
func (s *ProfileService) GetProfile(userID uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		// Phân biệt giữa "not found" và "internal error"
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFoundResponse()
		}
		return nil, utils.ErrInternalServerResponse()
	}
	return user, nil
}

// Cập nhật profile
func (s *ProfileService) UpdateProfile(userID uuid.UUID, name, phone, address, dateOfBirth, gender string, isEmailNotify *bool) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		// Phân biệt giữa "not found" và "internal error"
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFoundResponse()
		}
		return nil, utils.ErrInternalServerResponse()
	}

	// Cập nhật thông tin
	user.Name = name
	user.Phone = phone
	user.Address = address
	user.Gender = models.Gender(gender)

	if dateOfBirth != "" {
		parsed, err := time.Parse("2006-01-02", dateOfBirth)
		if err == nil {
			user.DateOfBirth = &parsed
		}
	}

	if isEmailNotify != nil {
		user.IsEmailNotify = *isEmailNotify
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, utils.ErrInternalServerResponse()
	}

	return user, nil
}
