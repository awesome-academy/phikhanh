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
	userID, _ := ctx.Get("user_id")
	uuid, _ := uuid.Parse(userID.(string))

	user, err := c.service.GetProfile(uuid)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "User not found")
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
// @Failure      400  {object}  utils.APIResponse
// @Router       /profile [put]
func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	uuid, _ := uuid.Parse(userID.(string))

	var req userDto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Trả về map[field]error
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	dobStr := ""
	if req.DateOfBirth != nil {
		dobStr = *req.DateOfBirth
	}

	user, err := c.service.UpdateProfile(
		uuid, req.Name, req.Phone, req.Address,
		dobStr, req.Gender, req.IsEmailNotify,
	)

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
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
