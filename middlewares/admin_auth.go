package middlewares

import (
	"net/http"
	"phikhanh/models"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

// allowedAdminRoles - Roles được phép truy cập admin portal
var allowedAdminRoles = map[string]bool{
	string(models.RoleStaff):   true,
	string(models.RoleManager): true,
	string(models.RoleAdmin):   true,
}

// AdminAuthMiddleware - Middleware kiểm tra JWT token và role
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("admin_token")
		if err != nil || token == "" {
			ctx.Redirect(http.StatusFound, "/admin/login")
			ctx.Abort()
			return
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			clearAdminCookies(ctx)
			ctx.Redirect(http.StatusFound, "/admin/login?error=Session+expired")
			ctx.Abort()
			return
		}

		if !allowedAdminRoles[claims.Role] {
			clearAdminCookies(ctx)
			ctx.Redirect(http.StatusFound, "/admin/login?error=Access+denied")
			ctx.Abort()
			return
		}

		ctx.Set("admin_id", claims.UserID)
		ctx.Set("admin_role", claims.Role)
		ctx.Next()
	}
}

// clearAdminCookies - Xóa tất cả admin cookies với HttpOnly=true (match với lúc set)
func clearAdminCookies(ctx *gin.Context) {
	ctx.SetCookie("admin_token", "", -1, "/admin", "", false, true)
	ctx.SetCookie("admin_name", "", -1, "/admin", "", false, true)
	ctx.SetCookie("admin_role", "", -1, "/admin", "", false, true)
}
