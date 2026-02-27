package admin

import (
	"log"
	"phikhanh/models"
	adminRepo "phikhanh/repositories/admin"
	"phikhanh/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo *adminRepo.AdminRepository
}

func NewAuthService(repo *adminRepo.AdminRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Login - Xác thực admin bằng email + password
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", utils.NewUnauthorizedError(utils.MsgInvalidCredentials)
		}
		return nil, "", utils.NewInternalServerError(err)
	}

	// Kiểm tra password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", utils.NewUnauthorizedError(utils.MsgInvalidCredentials)
	}

	// Tạo JWT token
	token, err := utils.GenerateToken(user.ID.String(), string(user.Role))
	if err != nil {
		log.Printf("[Admin Service] Token generation failed for '%s': %v", email, err)
		return nil, "", utils.NewInternalServerError(err)
	}

	return user, token, nil
}
