package user

import (
	"time"

	"phikhanh/models"
	userRepo "phikhanh/repositories/user"
	"phikhanh/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo *userRepo.AuthRepository
}

func NewAuthService(repo *userRepo.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Đăng ký tài khoản mới
func (s *AuthService) Register(citizenID, password, name, email, phone, address, dateOfBirth, gender string) (*models.User, error) {
	// Kiểm tra citizen_id đã tồn tại
	exists, err := s.repo.IsCitizenIDExists(citizenID)
	if err != nil {
		return nil, utils.ErrInternalServerResponse()
	}
	if exists {
		return nil, utils.ErrCitizenIDExistsResponse()
	}

	// Kiểm tra email đã tồn tại
	exists, err = s.repo.IsEmailExists(email)
	if err != nil {
		return nil, utils.ErrInternalServerResponse()
	}
	if exists {
		return nil, utils.ErrEmailExistsResponse()
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, utils.ErrInternalServerResponse()
	}

	// Parse date of birth
	var dob *time.Time
	if dateOfBirth != "" {
		parsed, err := time.Parse("2006-01-02", dateOfBirth)
		if err == nil {
			dob = &parsed
		}
	}

	// Tạo user mới
	user := &models.User{
		CitizenID:     citizenID,
		PasswordHash:  string(hashedPassword),
		Name:          name,
		Email:         email,
		Phone:         phone,
		Address:       address,
		DateOfBirth:   dob,
		Gender:        models.Gender(gender),
		Role:          models.RoleCitizen,
		IsEmailNotify: true,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, utils.ErrInternalServerResponse()
	}

	return user, nil
}

// Đăng nhập
func (s *AuthService) Login(citizenID, password string) (*models.User, string, error) {
	// Tìm user
	user, err := s.repo.FindByCitizenID(citizenID)
	if err != nil {
		// Phân biệt giữa "not found" và "internal error"
		if err == gorm.ErrRecordNotFound {
			return nil, "", utils.ErrInvalidCredentialsResponse()
		}
		return nil, "", utils.ErrInternalServerResponse()
	}

	// Kiểm tra password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", utils.ErrInvalidCredentialsResponse()
	}

	// Tạo JWT token
	token, err := utils.GenerateToken(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, "", utils.ErrInternalServerResponse()
	}

	return user, token, nil
}
