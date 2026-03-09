package admin

import (
	"log"
	"math"
	adminDto "phikhanh/dto/admin"
	"phikhanh/models"
	adminRepo "phikhanh/repositories/admin"
	"phikhanh/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const UserPageSize = 10

type UserService struct {
	repo         *adminRepo.UserRepository
	emailService *utils.EmailService
}

func NewUserService(repo *adminRepo.UserRepository) *UserService {
	return &UserService{
		repo:         repo,
		emailService: utils.NewEmailService(),
	}
}

// GetList - Lấy danh sách users với filter và pagination
func (s *UserService) GetList(role string, page int) (*adminDto.UserListResult, error) {
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * UserPageSize

	users, total, err := s.repo.FindAllWithFilter(role, offset, UserPageSize)
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(UserPageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	items := make([]adminDto.UserListItem, 0, len(users))
	for _, user := range users {
		items = append(items, adminDto.UserListItem{
			ID:        user.ID.String(),
			CitizenID: user.CitizenID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      string(user.Role),
			Status:    "Active",
		})
	}

	return &adminDto.UserListResult{
		Items:       items,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		Role:        role,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}, nil
}

// GetDetail - Lấy chi tiết user
func (s *UserService) GetDetail(id string) (*adminDto.UserDetail, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFoundError("User not found")
		}
		return nil, utils.NewInternalServerError(err)
	}

	detail := &adminDto.UserDetail{
		ID:        user.ID.String(),
		CitizenID: user.CitizenID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		Role:      string(user.Role),
		Gender:    string(user.Gender),
		CreatedAt: user.CreatedAt.Format(time.DateTime),
		UpdatedAt: user.UpdatedAt.Format(time.DateTime),
	}

	if user.DateOfBirth != nil {
		detail.DateOfBirth = user.DateOfBirth.Format("2006-01-02")
	}

	if user.Department != nil {
		detail.DepartmentName = user.Department.Name
	}

	if user.DepartmentID != nil {
		detail.DepartmentID = user.DepartmentID.String()
	}

	return detail, nil
}

// Create - Tạo user mới
func (s *UserService) Create(user *models.User, rawPassword string) error {
	// Validate required fields
	if user.CitizenID == "" {
		return utils.NewBadRequestError("Citizen ID is required")
	}
	if user.Name == "" {
		return utils.NewBadRequestError("Name is required")
	}
	if user.Email == "" {
		return utils.NewBadRequestError("Email is required")
	}

	// Kiểm tra citizen_id đã tồn tại
	exists, err := s.repo.IsCitizenIDExists(user.CitizenID)
	if err != nil {
		return utils.NewInternalServerError(err)
	}
	if exists {
		return utils.NewBadRequestError("Citizen ID already exists")
	}

	// Kiểm tra email đã tồn tại
	exists, err = s.repo.IsEmailExists(user.Email)
	if err != nil {
		return utils.NewInternalServerError(err)
	}
	if exists {
		return utils.NewBadRequestError("Email already exists")
	}

	// For staff role, department is required
	if user.Role == models.RoleStaff && user.DepartmentID == nil {
		return utils.NewBadRequestError("Department is required for staff")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewInternalServerError(err)
	}
	user.PasswordHash = string(hashedPassword)

	if err := s.repo.Create(user); err != nil {
		return utils.NewInternalServerError(err)
	}

	// Send welcome email async
	go s.sendWelcomeEmailAsync(user, rawPassword)

	return nil
}

// sendWelcomeEmailAsync - Gửi welcome email sau khi tạo user
func (s *UserService) sendWelcomeEmailAsync(user *models.User, rawPassword string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Welcome Email] Panic recovered: %v", r)
		}
	}()

	if user.Email == "" {
		log.Printf("[Welcome Email] Skipping: empty email for user %s", user.CitizenID)
		return
	}

	err := s.emailService.SendWelcomeEmail(user, rawPassword)
	if err != nil {
		errorHandler := utils.GetEmailErrorHandler()
		errorHandler.ReportError(user.CitizenID, user.Email, err)
	}
}

// Update - Cập nhật user
func (s *UserService) Update(user *models.User) error {
	// Validate required fields
	if user.CitizenID == "" {
		return utils.NewBadRequestError("Citizen ID is required")
	}
	if user.Name == "" {
		return utils.NewBadRequestError("Name is required")
	}
	if user.Email == "" {
		return utils.NewBadRequestError("Email is required")
	}

	// Kiểm tra citizen_id đã tồn tại (exclude current user)
	exists, err := s.repo.IsCitizenIDExistsExcept(user.CitizenID, user.ID.String())
	if err != nil {
		return utils.NewInternalServerError(err)
	}
	if exists {
		return utils.NewBadRequestError("Citizen ID already exists")
	}

	// Kiểm tra email đã tồn tại (exclude current user)
	exists, err = s.repo.IsEmailExistsExcept(user.Email, user.ID.String())
	if err != nil {
		return utils.NewInternalServerError(err)
	}
	if exists {
		return utils.NewBadRequestError("Email already exists")
	}

	// For staff role, department is required
	if user.Role == models.RoleStaff && user.DepartmentID == nil {
		return utils.NewBadRequestError("Department is required for staff")
	}

	// Prevent changing admin to other roles
	currentUser, err := s.repo.FindByID(user.ID.String())
	if err != nil {
		return utils.NewInternalServerError(err)
	}

	if currentUser.Role == models.RoleAdmin && user.Role != models.RoleAdmin {
		return utils.NewBadRequestError("Cannot change admin role to other roles")
	}

	if err := s.repo.Update(user); err != nil {
		return utils.NewInternalServerError(err)
	}

	return nil
}

// Delete - Soft delete user
func (s *UserService) Delete(id string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewNotFoundError("User not found")
		}
		return utils.NewInternalServerError(err)
	}

	// Prevent deleting admin users
	if user.Role == models.RoleAdmin {
		return utils.NewBadRequestError("Cannot delete admin users")
	}

	if err := s.repo.SoftDelete(id); err != nil {
		return utils.NewInternalServerError(err)
	}

	return nil
}

// GetAvailableRoles - Lấy danh sách roles (bao gồm citizen)
func (s *UserService) GetAvailableRoles() []adminDto.RoleOption {
	return []adminDto.RoleOption{
		{Value: "citizen", Label: "Citizen"},
		{Value: "staff", Label: "Staff"},
		{Value: "manager", Label: "Manager"},
		{Value: "admin", Label: "Admin"},
	}
}

// GetAvailableGenders - Lấy danh sách genders
func (s *UserService) GetAvailableGenders() []adminDto.GenderOption {
	return []adminDto.GenderOption{
		{Value: "male", Label: "Male"},
		{Value: "female", Label: "Female"},
		{Value: "other", Label: "Other"},
	}
}

// GetDepartments - Lấy danh sách departments
func (s *UserService) GetDepartments() ([]adminDto.DepartmentOption, error) {
	departments, err := s.repo.GetDepartments()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	result := make([]adminDto.DepartmentOption, 0, len(departments))
	for _, dept := range departments {
		result = append(result, adminDto.DepartmentOption{
			ID:   dept.ID.String(),
			Name: dept.Name,
		})
	}

	return result, nil
}

// ValidateAndPrepareForUpdate - Validate và prepare user data cho update
// Enforce role-based restrictions
func (s *UserService) ValidateAndPrepareForUpdate(user *models.User, currentUserRole string) error {
	// Validate required fields
	if user.CitizenID == "" {
		return utils.NewBadRequestError("Citizen ID is required")
	}
	if user.Name == "" {
		return utils.NewBadRequestError("Name is required")
	}
	if user.Email == "" {
		return utils.NewBadRequestError("Email is required")
	}

	// Get current user to check existing data
	currentUser, err := s.repo.FindByID(user.ID.String())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewNotFoundError("User not found")
		}
		return utils.NewInternalServerError(err)
	}

	// ROLE RESTRICTION: Only admin can change role
	if currentUserRole != "admin" {
		// Non-admin cannot change role
		user.Role = currentUser.Role
	}

	// ROLE RESTRICTION: Only admin/manager can set department
	if currentUserRole != "admin" && currentUserRole != "manager" {
		user.DepartmentID = currentUser.DepartmentID
	}

	// VALIDATION: For staff role, department is required
	if user.Role == models.RoleStaff && user.DepartmentID == nil {
		return utils.NewBadRequestError("Department is required for staff role")
	}

	// Kiểm tra citizen_id đã tồn tại (exclude current user)
	exists, err := s.repo.IsCitizenIDExistsExcept(user.CitizenID, user.ID.String())
	if err != nil {
		return utils.NewInternalServerError(err)
	}
	if exists {
		return utils.NewBadRequestError("Citizen ID already exists")
	}

	// Kiểm tra email đã tồn tại (exclude current user)
	exists, err = s.repo.IsEmailExistsExcept(user.Email, user.ID.String())
	if err != nil {
		return utils.NewInternalServerError(err)
	}
	if exists {
		return utils.NewBadRequestError("Email already exists")
	}

	// Prevent changing admin to other roles
	if currentUser.Role == models.RoleAdmin && user.Role != models.RoleAdmin {
		return utils.NewBadRequestError("Cannot change admin role to other roles")
	}

	return nil
}

// CanEditRole - Check if current user can edit role field
func (s *UserService) CanEditRole(currentUserRole string) bool {
	return currentUserRole == "admin"
}

// CanEditDepartment - Check if current user can edit department field
func (s *UserService) CanEditDepartment(currentUserRole string) bool {
	return currentUserRole == "admin" || currentUserRole == "manager"
}
