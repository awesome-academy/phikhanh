package user

import (
	"errors"
	"time"

	"phikhanh/models"
	userRepo "phikhanh/repositories/user"
	"phikhanh/utils"

	"golang.org/x/crypto/bcrypt"
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
		return nil, err
	}
	if exists {
		return nil, errors.New("citizen_id already exists")
	}

	// Kiểm tra email đã tồn tại
	exists, err = s.repo.IsEmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return user, nil
}

// Đăng nhập
func (s *AuthService) Login(citizenID, password string) (*models.User, string, error) {
	// Tìm user
	user, err := s.repo.FindByCitizenID(citizenID)
	if err != nil {
		return nil, "", errors.New("invalid citizen_id or password")
	}

	// Kiểm tra password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid citizen_id or password")
	}

	// Tạo JWT token
	token, err := utils.GenerateToken(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
