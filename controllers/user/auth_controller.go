package user

import (
	"net/http"

	userDto "phikhanh/dto/user"
	userSvc "phikhanh/services/user"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *userSvc.AuthService
}

func NewAuthController(service *userSvc.AuthService) *AuthController {
	return &AuthController{service: service}
}

// Register - Đăng ký tài khoản
// POST /api/v1/auth/register
func (c *AuthController) Register(ctx *gin.Context) {
	var req userDto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Trả về map[field]error
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	user, err := c.service.Register(
		req.CitizenID, req.Password, req.Name, req.Email,
		req.Phone, req.Address, req.DateOfBirth, req.Gender,
	)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response := userDto.UserInfo{
		ID:        user.ID.String(),
		CitizenID: user.CitizenID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      string(user.Role),
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Register successful", response)
}

// Login - Đăng nhập
// POST /api/v1/auth/login
func (c *AuthController) Login(ctx *gin.Context) {
	var req userDto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Trả về map[field]error
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	user, token, err := c.service.Login(req.CitizenID, req.Password)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	response := userDto.LoginResponse{
		Token: token,
		User: userDto.UserInfo{
			ID:        user.ID.String(),
			CitizenID: user.CitizenID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Role:      string(user.Role),
		},
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Login successful", response)
}

// Logout - Đăng xuất
// POST /api/v1/auth/logout
func (c *AuthController) Logout(ctx *gin.Context) {
	// JWT stateless, client tự xóa token
	utils.SuccessResponse(ctx, http.StatusOK, "Logout successful", nil)
}
