package admin

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"phikhanh/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	redirectDepartments = "/admin/departments"
	redirectServices    = "/admin/services"
	redirectUsers       = "/admin/users"
)

// formatErrorMessage - Format error thành human-readable message cho admin SSR
func formatErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		if appErr.StatusCode >= 400 && appErr.StatusCode < 500 {
			return appErr.Message
		}
		return "An error occurred while processing your request"
	}

	// Handle validation errors từ gin binding
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		messages := make([]string, 0, len(valErrs))
		for _, msg := range utils.FormatValidationErrorsMap(err) {
			messages = append(messages, msg)
		}
		return strings.Join(messages, "; ")
	}

	return err.Error()
}

// setFlashError - Redirect với error message dùng query string
func setFlashError(ctx *gin.Context, message string, redirectTo string) {
	ctx.Redirect(http.StatusFound, redirectTo+"?error="+url.QueryEscape(message))
}

// setFlashSuccess - Redirect với success message dùng query string
func setFlashSuccess(ctx *gin.Context, message string, redirectTo string) {
	ctx.Redirect(http.StatusFound, redirectTo+"?success="+url.QueryEscape(message))
}

// getCsrfToken - Get CSRF token từ context
func getCsrfToken(ctx *gin.Context) string {
	token, _ := ctx.Get("csrf_token")
	if token == nil {
		return ""
	}
	return token.(string)
}

// parseUUID - Parse UUID từ path param, redirect với error nếu không hợp lệ
func parseUUID(ctx *gin.Context, redirectTo string) (uuid.UUID, bool) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		setFlashError(ctx, "Invalid ID format", redirectTo)
		return uuid.Nil, false
	}
	return id, true
}

// parseDepartmentID - Parse UUID cho department routes
func parseDepartmentID(ctx *gin.Context) (uuid.UUID, bool) {
	return parseUUID(ctx, redirectDepartments)
}

// parseServiceID - Parse UUID cho service routes
func parseServiceID(ctx *gin.Context) (uuid.UUID, bool) {
	return parseUUID(ctx, redirectServices)
}
