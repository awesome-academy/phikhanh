package user

import (
	"net/http"

	userDto "phikhanh/dto/user"
	userSvc "phikhanh/services/user"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplicationController struct {
	service *userSvc.ApplicationService
}

func NewApplicationController(service *userSvc.ApplicationService) *ApplicationController {
	return &ApplicationController{service: service}
}

// SubmitApplication godoc
// @Summary      Nộp hồ sơ
// @Description  Công dân nộp hồ sơ dịch vụ công
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body userDto.SubmitAppRequest true "Thông tin hồ sơ"
// @Success      201  {object}  utils.APIResponse{data=userDto.SubmitAppResponse}
// @Failure      400  {object}  utils.ValidationErrors
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /applications [post]
func (c *ApplicationController) SubmitApplication(ctx *gin.Context) {
	// Lấy user_id từ context (set bởi AuthMiddleware)
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		svcErr := utils.NewUnauthorizedError("Unauthorized")
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		svcErr := utils.NewUnauthorizedError("Invalid user ID")
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		svcErr := utils.NewBadRequestError("Invalid user ID format")
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	// Parse và validate request body
	var req userDto.SubmitAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	// Gọi service xử lý business logic
	response, err := c.service.SubmitApplication(req, userID)
	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			return
		}
		svcErr := utils.NewInternalServerError(err)
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Application submitted successfully", response)
}
