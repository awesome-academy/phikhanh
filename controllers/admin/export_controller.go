package admin

import (
	"encoding/csv"
	"net/http"

	adminSvc "phikhanh/services/admin"

	"github.com/gin-gonic/gin"
)

// utf8BOM - Required để Excel hiển thị đúng tiếng Việt
var utf8BOM = []byte("\xEF\xBB\xBF")

type ExportController struct {
	service *adminSvc.ExportService
}

func NewExportController(service *adminSvc.ExportService) *ExportController {
	return &ExportController{service: service}
}

// writeCSV - Helper: set headers, write BOM, write CSV rows
func (c *ExportController) writeCSV(ctx *gin.Context, filename string, rows [][]string) {
	ctx.Header("Content-Type", "text/csv; charset=utf-8")
	ctx.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	ctx.Status(http.StatusOK)

	// Write UTF-8 BOM cho Excel
	ctx.Writer.Write(utf8BOM)

	w := csv.NewWriter(ctx.Writer)
	if err := w.WriteAll(rows); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	w.Flush()
}

// GET /admin/export/citizens
func (c *ExportController) ExportCitizens(ctx *gin.Context) {
	rows, err := c.service.ExportCitizens()
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/dashboard")
		return
	}
	c.writeCSV(ctx, "citizens.csv", rows)
}

// GET /admin/export/applications
func (c *ExportController) ExportApplications(ctx *gin.Context) {
	rows, err := c.service.ExportApplications()
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/dashboard")
		return
	}
	c.writeCSV(ctx, "applications.csv", rows)
}

// GET /admin/export/services
func (c *ExportController) ExportServices(ctx *gin.Context) {
	rows, err := c.service.ExportServices()
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/dashboard")
		return
	}
	c.writeCSV(ctx, "services.csv", rows)
}

// GET /admin/export/departments
func (c *ExportController) ExportDepartments(ctx *gin.Context) {
	rows, err := c.service.ExportDepartments()
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/dashboard")
		return
	}
	c.writeCSV(ctx, "departments.csv", rows)
}

// GET /admin/export/staffs
func (c *ExportController) ExportStaff(ctx *gin.Context) {
	rows, err := c.service.ExportStaff()
	if err != nil {
		setFlashError(ctx, formatErrorMessage(err), "/admin/dashboard")
		return
	}
	c.writeCSV(ctx, "staffs.csv", rows)
}
