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

// Register godoc
// @Summary      Đăng ký tài khoản
// @Description  Đăng ký tài khoản công dân mới
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body userDto.RegisterRequest true "Thông tin đăng ký"
// @Success      201  {object}  utils.APIResponse{data=userDto.UserInfo}
// @Failure      400  {object}  utils.ValidationErrors
// @Router       /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req userDto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	user, err := c.service.Register(
		req.CitizenID, req.Password, req.Name, req.Email,
		req.Phone, req.Address, req.DateOfBirth, req.Gender,
	)

	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			return
		}
		// Fallback for unexpected errors
		svcErr := utils.NewInternalServerError(err)
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
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

// Login godoc
// @Summary      Đăng nhập
// @Description  Đăng nhập bằng CCCD và mật khẩu
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body userDto.LoginRequest true "Thông tin đăng nhập"
// @Success      200  {object}  utils.APIResponse{data=userDto.LoginResponse}
// @Failure      400  {object}  utils.ValidationErrors
// @Failure      401  {object}  utils.APIResponse
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req userDto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	user, token, err := c.service.Login(req.CitizenID, req.Password)
	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			return
		}
		// Fallback for unexpected errors
		svcErr := utils.NewInternalServerError(err)
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
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

// Logout godoc
// @Summary      Đăng xuất
// @Description  Đăng xuất khỏi hệ thống (client xóa token)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Router       /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// JWT stateless, client tự xóa token
	utils.SuccessResponse(ctx, http.StatusOK, "Logout successful", nil)
}
