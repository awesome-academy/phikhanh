package user

import (
	"net/http"
	"strconv"

	"phikhanh/services"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationController struct {
	service *services.NotificationService
}

func NewNotificationController(service *services.NotificationService) *NotificationController {
	return &NotificationController{service: service}
}

// GetNotifications godoc
// @Summary      Danh sách thông báo
// @Tags         Notifications
// @Security     BearerAuth
// @Param        page query int false "Page" default(1)
// @Success      200  {object}  utils.APIResponse
// @Router       /notifications [get]
func (c *NotificationController) GetNotifications(ctx *gin.Context) {
	userID, svcErr := utils.ExtractUserID(ctx)
	if svcErr != nil {
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil && page <= 0 {
		page = 1
	}

	result, err := c.service.GetNotifications(userID, page)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(ctx, appErr.StatusCode, appErr.Message)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get notifications")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Get notifications successful", result)
}

// MarkAsRead godoc
// @Summary      Đánh dấu đã đọc
// @Tags         Notifications
// @Security     BearerAuth
// @Param        id path string true "Notification ID"
// @Success      200  {object}  utils.APIResponse
// @Router       /notifications/{id}/read [patch]
func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	userID, svcErr := utils.ExtractUserID(ctx)
	if svcErr != nil {
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	notifID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := c.service.MarkAsRead(notifID, userID); err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Notification not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Marked as read", nil)
}

// MarkAllAsRead godoc
// @Summary      Đánh dấu tất cả đã đọc
// @Tags         Notifications
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse
// @Router       /notifications/read-all [patch]
func (c *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID, svcErr := utils.ExtractUserID(ctx)
	if svcErr != nil {
		utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
		return
	}

	if err := c.service.MarkAllAsRead(userID); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to mark all as read")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "All notifications marked as read", nil)
}
