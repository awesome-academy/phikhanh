package middlewares

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// RequireRole - Middleware kiểm tra role của user
// Nếu user không có required role, redirect về dashboard với error message
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy role từ context (được set bởi auth middleware)
		role, exists := ctx.Get("admin_role")
		if !exists {
			ctx.Redirect(http.StatusFound, "/admin/dashboard?error="+url.QueryEscape("Unauthorized access"))
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			ctx.Redirect(http.StatusFound, "/admin/dashboard?error="+url.QueryEscape("Invalid role format"))
			return
		}

		// Kiểm tra role có trong danh sách allowed roles
		allowed := false
		for _, ar := range allowedRoles {
			if roleStr == ar {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.Redirect(http.StatusFound, "/admin/dashboard?error="+url.QueryEscape("Access denied. Insufficient permissions"))
			return
		}

		ctx.Next()
	}
}
