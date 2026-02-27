package admin

import (
	"net/http"
	"net/url"

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
	// Đã login → redirect dashboard
	if token, err := ctx.Cookie("admin_token"); err == nil && token != "" {
		if _, err := utils.VerifyToken(token); err == nil {
			ctx.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}
	}

	ctx.HTML(http.StatusOK, "admin/auth/login.html", gin.H{
		"Error": ctx.Query("error"),
	})
}

// ProcessLogin - Xử lý login admin
func (c *AuthController) ProcessLogin(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if email == "" || password == "" {
		redirectWithError(ctx, "Email and password are required")
		return
	}

	user, token, err := c.service.Login(email, password)
	if err != nil {
		if svcErr, ok := err.(*utils.ServiceError); ok {
			redirectWithError(ctx, svcErr.Message)
			return
		}
		redirectWithError(ctx, "An error occurred")
		return
	}

	ctx.SetCookie("admin_token", token, cookieMaxAge, cookiePath, "", false, true)
	ctx.SetCookie("admin_name", user.Name, cookieMaxAge, cookiePath, "", false, false)
	ctx.SetCookie("admin_role", string(user.Role), cookieMaxAge, cookiePath, "", false, false)

	ctx.Redirect(http.StatusFound, "/admin/dashboard")
}

// ProcessLogout - Xử lý logout admin
func (c *AuthController) ProcessLogout(ctx *gin.Context) {
	ctx.SetCookie("admin_token", "", -1, cookiePath, "", false, true)
	ctx.SetCookie("admin_name", "", -1, cookiePath, "", false, false)
	ctx.SetCookie("admin_role", "", -1, cookiePath, "", false, false)
	ctx.Redirect(http.StatusFound, "/admin/login")
}

// redirectWithError - Redirect về login với error message
func redirectWithError(ctx *gin.Context, message string) {
	ctx.Redirect(http.StatusFound, "/admin/login?error="+url.QueryEscape(message))
}
