package admin

import (
	"net/http"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type ApplicationController struct{}

func NewApplicationController() *ApplicationController {
	return &ApplicationController{}
}

func (c *ApplicationController) ShowApplications(ctx *gin.Context) {
	utils.RenderHTML(ctx, http.StatusOK, "admin/applications.html", utils.GetAdminData(ctx, "Applications", "applications"))
}
