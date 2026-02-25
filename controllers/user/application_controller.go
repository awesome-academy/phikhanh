package user

import (
	"net/http"

	userDto "phikhanh/dto/user"
	userSvc "phikhanh/services/user"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
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
	userID, svcErr := utils.ExtractUserID(ctx)
	if svcErr != nil {
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	var req userDto.SubmitAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

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

// GetMyApplications godoc
// @Summary      Danh sách hồ sơ của tôi
// @Description  Lấy danh sách hồ sơ đã nộp của công dân
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query int    false "Page number" default(1)
// @Param        limit  query int    false "Items per page" default(10)
// @Param        status query string false "Filter by status" Enums(Received,Processing,Supplement_Required,Approved,Rejected)
// @Success      200  {object}  utils.APIResponse{data=userDto.MyAppListResponse}
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /applications/me [get]
func (c *ApplicationController) GetMyApplications(ctx *gin.Context) {
	userID, svcErr := utils.ExtractUserID(ctx)
	if svcErr != nil {
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	var req userDto.MyAppListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	response, err := c.service.GetMyApplications(req, userID)
	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			return
		}
		svcErr := utils.NewInternalServerError(err)
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Get my applications successful", response)
}
