package admin

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	redirectDepartments = "/admin/departments"
	redirectServices    = "/admin/services"
)

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
