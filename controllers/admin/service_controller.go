package admin

import (
	"net/http"
	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

type ServiceController struct{}

func NewServiceController() *ServiceController {
	return &ServiceController{}
}

func (c *ServiceController) ShowServices(ctx *gin.Context) {
	utils.RenderHTML(ctx, http.StatusOK, "admin/services.html", utils.GetAdminData(ctx, "Services", "services"))
}
