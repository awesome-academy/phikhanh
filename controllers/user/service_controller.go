package user

import (
	"net/http"

	userDto "phikhanh/dto/user"
	userSvc "phikhanh/services/user"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceController struct {
	service *userSvc.ServiceService
}

func NewServiceController(service *userSvc.ServiceService) *ServiceController {
	return &ServiceController{service: service}
}

// GetServiceList godoc
// @Summary      Danh sách dịch vụ
// @Description  Lấy danh sách dịch vụ với pagination và filter
// @Tags         Services
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        keyword query string false "Search by name or code"
// @Param        sector query string false "Filter by sector"
// @Param        department_id query string false "Filter by department ID"
// @Success      200  {object}  utils.APIResponse{data=userDto.ServiceListResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /services [get]
func (c *ServiceController) GetServiceList(ctx *gin.Context) {
	var req userDto.ServiceListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(ctx, utils.FormatValidationErrorsMap(err))
		return
	}

	// Enforce pagination defaults BEFORE calling service
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	response, err := c.service.GetServiceList(req)
	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			return
		}
		svcErr := utils.NewInternalServerError(err)
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Get service list successful", response)
}

// GetServiceDetail godoc
// @Summary      Chi tiết dịch vụ
// @Description  Lấy thông tin chi tiết dịch vụ
// @Tags         Services
// @Accept       json
// @Produce      json
// @Param        id path string true "Service ID"
// @Success      200  {object}  utils.APIResponse{data=userDto.ServiceDetailResponse}
// @Failure      404  {object}  utils.APIResponse
// @Router       /services/{id} [get]
func (c *ServiceController) GetServiceDetail(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		svcErr := utils.NewBadRequestError("Invalid service ID")
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	response, err := c.service.GetServiceDetail(id)
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

	utils.SuccessResponse(ctx, http.StatusOK, "Get service detail successful", response)
}
