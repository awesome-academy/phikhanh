package user

import (
	"net/http"
	"time"

	userDto "phikhanh/dto/user"
	userSvc "phikhanh/services/user"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileController struct {
	service *userSvc.ProfileService
}

func NewProfileController(service *userSvc.ProfileService) *ProfileController {
	return &ProfileController{service: service}
}

// GetProfile godoc
// @Summary      Lấy thông tin hồ sơ
// @Description  Lấy thông tin hồ sơ công dân
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse{data=userDto.ProfileResponse}
// @Failure      401  {object}  utils.APIResponse
// @Router       /profile [get]
func (c *ProfileController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, 401, "Unauthorized")
		return
	}

	parsedUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid user ID")
		return
	}

	user, err := c.service.GetProfile(parsedUUID)
	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			return
		}
		svcErr := utils.ErrInternalServerResponse()
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	var dobStr *string
	if user.DateOfBirth != nil {
		formatted := user.DateOfBirth.Format("2006-01-02")
		dobStr = &formatted
	}

	response := userDto.ProfileResponse{
		ID:            user.ID.String(),
		CitizenID:     user.CitizenID,
		Name:          user.Name,
		Email:         user.Email,
		Phone:         user.Phone,
		Address:       user.Address,
		DateOfBirth:   dobStr,
		Gender:        string(user.Gender),
		Role:          string(user.Role),
		IsEmailNotify: user.IsEmailNotify,
		CreatedAt:     user.CreatedAt.Format(time.RFC3339),
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Get profile successful", response)
}

// UpdateProfile godoc
// @Summary      Cập nhật hồ sơ
// @Description  Cập nhật thông tin hồ sơ công dân
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body userDto.UpdateProfileRequest true "Thông tin cập nhật"
// @Success      200  {object}  utils.APIResponse{data=userDto.ProfileResponse}
// @Failure      400  {object}  utils.ValidationErrors
// @Failure      401  {object}  utils.APIResponse
// @Router       /profile [put]
func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, 401, "Unauthorized")
		return
	}

	parsedUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.ErrorResponse(ctx, 400, "Invalid user ID")
		return
	}

	var req userDto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	dobStr := ""
	if req.DateOfBirth != nil {
		dobStr = *req.DateOfBirth
	}

	addressStr := ""
	if req.Address != nil {
		addressStr = *req.Address
	}

	user, err := c.service.UpdateProfile(
		parsedUUID, req.Name, req.Phone, addressStr,
		dobStr, req.Gender, req.IsEmailNotify,
	)

	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			return
		}
		svcErr := utils.ErrInternalServerResponse()
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	var dobResponse *string
	if user.DateOfBirth != nil {
		formatted := user.DateOfBirth.Format("2006-01-02")
		dobResponse = &formatted
	}

	response := userDto.ProfileResponse{
		ID:            user.ID.String(),
		CitizenID:     user.CitizenID,
		Name:          user.Name,
		Email:         user.Email,
		Phone:         user.Phone,
		Address:       user.Address,
		DateOfBirth:   dobResponse,
		Gender:        string(user.Gender),
		Role:          string(user.Role),
		IsEmailNotify: user.IsEmailNotify,
		CreatedAt:     user.CreatedAt.Format(time.RFC3339),
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Update profile successful", response)
}
