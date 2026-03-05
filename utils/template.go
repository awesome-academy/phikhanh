package utils

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// statusBadgeClass - Map status -> CSS classes
var statusBadgeClassMap = map[string]string{
	"Received":            "bg-gray-100 text-gray-700",
	"Processing":          "bg-yellow-100 text-yellow-700",
	"Supplement_Required": "bg-orange-100 text-orange-700",
	"Approved":            "bg-green-100 text-green-700",
	"Rejected":            "bg-red-100 text-red-700",
}

var templateFuncs = template.FuncMap{
	// codeNumber - Extract số từ code format PREFIX-XXX → XXX
	"codeNumber": func(code string) string {
		if idx := strings.LastIndex(code, "-"); idx != -1 {
			return code[idx+1:]
		}
		return code
	},
	// add - Cộng hai số nguyên (dùng cho pagination)
	"add": func(a, b int) int { return a + b },
	// sub - Trừ hai số nguyên (dùng cho pagination)
	"sub": func(a, b int) int { return a - b },
	// statusBadge - Lấy CSS classes cho status badge
	"statusBadge": func(status string) string {
		if class, ok := statusBadgeClassMap[status]; ok {
			return class
		}
		return "bg-gray-100 text-gray-700"
	},
}

var pageTemplates = map[string]*template.Template{}

// LoadTemplates - Load tất cả templates, mỗi page có template set riêng
func LoadTemplates(dir string) *template.Template {
	layoutFile := filepath.Join(dir, "admin", "layout", "base.html")
	loginFile := filepath.Join(dir, "admin", "auth", "login.html")

	// Fail fast nếu file bắt buộc không tồn tại
	for _, f := range []string{layoutFile, loginFile} {
		if _, err := os.Stat(f); err != nil {
			log.Fatalf("Required template file not found: %s", f)
		}
	}

	// Parse login standalone với funcs
	loginContent, err := os.ReadFile(loginFile)
	if err != nil {
		log.Fatalf("Failed to read login.html: %v", err)
	}
	loginTmpl, err := template.New("admin/auth/login.html").Funcs(templateFuncs).Parse(string(loginContent))
	if err != nil {
		log.Fatalf("Failed to parse login.html: %v", err)
	}
	pageTemplates["admin/auth/login.html"] = loginTmpl
	log.Printf("✓ Template loaded: admin/auth/login.html")

	// Collect tất cả page files (bỏ qua layout, login, và email templates)
	var pageFiles []string
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return walkErr
		}

		// Skip layout và login files
		if path == layoutFile || path == loginFile {
			return nil
		}

		// SKIP email templates (không cần {{ define "content" }} block)
		if strings.Contains(path, "email/") {
			return nil
		}

		pageFiles = append(pageFiles, path)
		return nil
	})

	// Mỗi page parse riêng biệt cùng layout
	for _, pageFile := range pageFiles {
		name, _ := filepath.Rel(dir, pageFile)

		// Thêm Funcs vào template set
		tmpl, err := template.New("base").Funcs(templateFuncs).ParseFiles(layoutFile, pageFile)
		if err != nil {
			log.Fatalf("Failed to parse template [%s]: %v", name, err)
		}

		if tmpl.Lookup("content") == nil {
			log.Fatalf("Template %q missing {{ define \"content\" }}...{{ end }} block", name)
		}

		pageTemplates[name] = tmpl
		log.Printf("✓ Template loaded: %s", name)
	}

	return template.New("root")
}

// GetTemplate - Lấy template set theo tên page
func GetTemplate(name string) (*template.Template, error) {
	tmpl, ok := pageTemplates[name]
	if !ok {
		return nil, fmt.Errorf("template %q not found", name)
	}
	return tmpl, nil
}
