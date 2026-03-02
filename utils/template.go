package utils

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// pageTemplates - Lưu từng page template set riêng biệt
// key: template name (e.g. "admin/dashboard.html")
// value: template set chứa cả "base" và "content"
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

	// Collect page files
	var pageFiles []string
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return walkErr
		}
		if path == layoutFile || path == loginFile {
			return nil
		}
		pageFiles = append(pageFiles, path)
		return nil
	})

	// Parse mỗi page cùng layout vào template set riêng
	for _, pageFile := range pageFiles {
		name, _ := filepath.Rel(dir, pageFile)

		tmpl, err := template.New("base").ParseFiles(layoutFile, pageFile)
		if err != nil {
			log.Fatalf("Failed to parse template [%s]: %v", name, err)
		}

		// Kiểm tra "content" block tồn tại - fail fast
		if tmpl.Lookup("content") == nil {
			log.Fatalf("Template %q missing {{ define \"content\" }}...{{ end }} block", name)
		}

		// Lưu toàn bộ template set (bao gồm "base" + "content")
		pageTemplates[name] = tmpl
		log.Printf("✓ Template loaded: %s", name)
	}

	// Parse login standalone
	loginContent, err := os.ReadFile(loginFile)
	if err != nil {
		log.Fatalf("Failed to read login.html: %v", err)
	}
	loginTmpl := template.New("admin/auth/login.html")
	if _, err := loginTmpl.Parse(string(loginContent)); err != nil {
		log.Fatalf("Failed to parse login.html: %v", err)
	}
	pageTemplates["admin/auth/login.html"] = loginTmpl
	log.Printf("✓ Template loaded: admin/auth/login.html")

	// Trả về dummy root (Gin sẽ dùng GetTemplate thay thế)
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
