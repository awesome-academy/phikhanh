package utils

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	layoutDirective = `{{ template "admin/layout/base.html" . }}`
	contentBlock    = `{{ block "content" . }}{{ end }}`
	contentDefine   = `{{ define "content" }}`
	contentEnd      = `{{ end }}`
)

// LoadTemplates - Load templates với layout inheritance
// Mỗi page template được inline vào layout thành HTML hoàn chỉnh
func LoadTemplates(dir string) *template.Template {
	root := template.New("root")

	layoutContent := loadLayoutFile(dir)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return walkErr
		}

		name, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// Skip layout files
		if strings.Contains(name, "layout/") {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(raw)
		var finalHTML string

		if strings.Contains(content, layoutDirective) && layoutContent != "" {
			// Page dùng layout → inline content vào layout
			pageBody := extractContentBlock(content)
			finalHTML = strings.Replace(layoutContent, contentBlock, pageBody, 1)
		} else {
			// Standalone template (login.html)
			finalHTML = content
		}

		if _, parseErr := root.New(name).Parse(finalHTML); parseErr != nil {
			log.Fatalf("Template parse error [%s]: %v", name, parseErr)
		}

		log.Printf("✓ Template loaded: %s", name)
		return nil
	})

	if err != nil {
		log.Fatalf("Template walk error: %v", err)
	}

	return root
}

// loadLayoutFile - Đọc base layout file
func loadLayoutFile(dir string) string {
	layoutFile := filepath.Join(dir, "admin", "layout", "base.html")
	content, err := os.ReadFile(layoutFile)
	if err != nil {
		log.Printf("⚠ Layout file not found: %s", layoutFile)
		return ""
	}
	log.Printf("✓ Layout loaded: admin/layout/base.html")
	return string(content)
}

// extractContentBlock - Lấy nội dung giữa {{ define "content" }} và {{ end }}
func extractContentBlock(content string) string {
	start := strings.Index(content, contentDefine)
	if start == -1 {
		return ""
	}

	body := content[start+len(contentDefine):]

	end := strings.LastIndex(body, contentEnd)
	if end == -1 {
		return body
	}

	return strings.TrimSpace(body[:end])
}
