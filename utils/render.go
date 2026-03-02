package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderHTML - Render template với data, thay thế ctx.HTML()
func RenderHTML(ctx *gin.Context, status int, templateName string, data any) {
	tmpl, err := GetTemplate(templateName)
	if err != nil {
		log.Printf("Template not found: %s", templateName)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(status)
	ctx.Header("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(ctx.Writer, data); err != nil {
		log.Printf("Template render error [%s]: %v", templateName, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
