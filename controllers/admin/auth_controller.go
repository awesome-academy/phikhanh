package admin

import (
	"net/http"

	adminDto "phikhanh/dto/admin"
	adminSvc "phikhanh/services/admin"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

const (
	cookieMaxAge = 3600 // 1 hour
	cookiePath   = "/admin"
)

type AuthController struct {
	service *adminSvc.AuthService
}

func NewAuthController(service *adminSvc.AuthService) *AuthController {
	return &AuthController{service: service}
}

// ShowLogin - Hiển thị trang login admin
func (c *AuthController) ShowLogin(ctx *gin.Context) {
	if token, err := ctx.Cookie("admin_token"); err == nil && token != "" {
		if _, err := utils.VerifyToken(token); err == nil {
			ctx.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}
	}

	utils.RenderHTML(ctx, http.StatusOK, "admin/auth/login.html", gin.H{
		"Error": ctx.Query("error"),
	})
}

// ProcessLogin - Xử lý login admin
func (c *AuthController) ProcessLogin(ctx *gin.Context) {
	var req adminDto.LoginRequest

	if err := ctx.ShouldBind(&req); err != nil {
		redirectWithError(ctx, formatErrorMessage(err))
		return
	}

	user, token, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		redirectWithError(ctx, formatErrorMessage(err))
		return
	}

	// Tất cả cookies đều HttpOnly=true vì chỉ dùng server-side cho SSR
	ctx.SetCookie("admin_token", token, cookieMaxAge, cookiePath, "", false, true)
	ctx.SetCookie("admin_name", user.Name, cookieMaxAge, cookiePath, "", false, true)
	ctx.SetCookie("admin_role", string(user.Role), cookieMaxAge, cookiePath, "", false, true)

	ctx.Redirect(http.StatusFound, "/admin/dashboard")
}

// ProcessLogout - Xử lý logout admin
func (c *AuthController) ProcessLogout(ctx *gin.Context) {
	// Xóa tất cả cookies với cùng HttpOnly=true
	ctx.SetCookie("admin_token", "", -1, cookiePath, "", false, true)
	ctx.SetCookie("admin_name", "", -1, cookiePath, "", false, true)
	ctx.SetCookie("admin_role", "", -1, cookiePath, "", false, true)
	ctx.Redirect(http.StatusFound, "/admin/login")
}

// redirectWithError - Redirect về login với error message
func redirectWithError(ctx *gin.Context, message string) {
	setFlashError(ctx, message, "/admin/login")
}
