package middlewares

import (
	"strings"

	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

// Middleware xác thực JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			svcErr := utils.NewUnauthorizedError("Authorization header required")
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			svcErr := utils.NewUnauthorizedError("Invalid authorization header format")
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			ctx.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.VerifyToken(token)
		if err != nil {
			svcErr := utils.NewUnauthorizedError("Invalid or expired token")
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			ctx.Abort()
			return
		}

		// Lưu thông tin user vào context (đảm bảo luôn tồn tại)
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
