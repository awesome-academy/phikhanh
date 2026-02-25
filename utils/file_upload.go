package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxFileSize = 10 << 20 // 10MB
	AllowedExts = ".jpg,.jpeg,.png,.pdf"
)

// getUploadDir - Lấy upload directory từ env hoặc dùng absolute path mặc định
func getUploadDir() string {
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		return dir
	}

	// Lấy absolute path từ working directory
	execDir, err := os.Getwd()
	if err != nil {
		return "./assets/images"
	}

	return filepath.Join(execDir, "assets", "images")
}

type UploadResult struct {
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	FileURL      string `json:"file_url"`
	FileSize     int64  `json:"file_size"`
	OriginalName string `json:"original_name"`
}

// UploadFile - Upload file và trả về đường dẫn
func UploadFile(file *multipart.FileHeader) (*UploadResult, error) {
	uploadDir := getUploadDir()

	// Validate file size
	if file.Size > MaxFileSize {
		return nil, NewBadRequestError("File size exceeds maximum allowed (10MB)")
	}

	// Sanitize filename trước khi validate extension
	sanitizedName := filepath.Base(file.Filename)
	sanitizedName = strings.ReplaceAll(sanitizedName, "\x00", "")

	// Validate file extension từ sanitized filename
	ext := strings.ToLower(filepath.Ext(sanitizedName))
	if !isAllowedExtension(ext) {
		return nil, NewBadRequestError(fmt.Sprintf("File type not allowed. Allowed types: %s", AllowedExts))
	}

	// Create upload directory
	if err := os.MkdirAll(uploadDir, 0700); err != nil {
		return nil, NewInternalServerError(err)
	}

	// Generate unique filename
	uniqueFileName := generateUniqueFileName(file.Filename)
	filePath := filepath.Join(uploadDir, uniqueFileName)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, NewInternalServerError(err)
	}
	defer src.Close()

	// Create destination file với restrictive permissions (chỉ owner đọc/ghi)
	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, NewInternalServerError(err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		// Đảm bảo đóng file trước khi xóa để tránh giữ handle mở trên một số OS
		_ = dst.Close()
		// Xóa file nếu copy thất bại
		_ = os.Remove(filePath)
		return nil, NewInternalServerError(err)
	}

	// Generate file URL
	baseURL := os.Getenv("FILE_BASE_URL")
	if baseURL == "" {
		baseURL = "/assets/images"
	}
	fileURL := fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), uniqueFileName)

	return &UploadResult{
		FileName:     uniqueFileName,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     file.Size,
		OriginalName: file.Filename,
	}, nil
}

// generateUniqueFileName - Tạo tên file unique với sanitized filename
func generateUniqueFileName(originalName string) string {
	// Lấy base filename để tránh path traversal (../../../etc/passwd)
	baseName := filepath.Base(originalName)

	// Loại bỏ null bytes để tránh null byte injection (file.jpg\x00.exe)
	baseName = strings.ReplaceAll(baseName, "\x00", "")

	// Extract extension từ sanitized filename
	ext := strings.ToLower(filepath.Ext(baseName))

	timestamp := time.Now().Format("20060102150405")
	uniqueID := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s%s", timestamp, uniqueID, ext)
}

// isAllowedExtension - Kiểm tra extension có được phép không
func isAllowedExtension(ext string) bool {
	allowedList := strings.Split(AllowedExts, ",")
	for _, allowed := range allowedList {
		if strings.TrimSpace(allowed) == ext {
			return true
		}
	}
	return false
}

// DeleteFile - Xóa file từ server với path validation
func DeleteFile(filePath string) error {
	uploadDir := getUploadDir()

	// Resolve symlinks cho allowed directory
	allowedDir, err := filepath.EvalSymlinks(uploadDir)
	if err != nil {
		// Nếu dir chưa tồn tại thì dùng Abs
		allowedDir, err = filepath.Abs(uploadDir)
		if err != nil {
			return NewInternalServerError(err)
		}
	}

	// Clean path trước
	cleanPath := filepath.Clean(filePath)

	// Resolve symlinks cho file path
	absPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		// File chưa tồn tại hoặc symlink broken
		if os.IsNotExist(err) {
			return NewNotFoundError("File not found")
		}
		absPath, err = filepath.Abs(cleanPath)
		if err != nil {
			return NewInternalServerError(err)
		}
	}

	// Dùng filepath.Rel để kiểm tra containment - an toàn hơn HasPrefix
	rel, err := filepath.Rel(allowedDir, absPath)
	if err != nil {
		return NewInternalServerError(err)
	}

	// Nếu relative path là ".." hoặc bắt đầu bằng ".." + path separator thì file nằm ngoài allowed directory
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return NewBadRequestError("Invalid file path")
	}

	// Kiểm tra file tồn tại
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return NewNotFoundError("File not found")
	}

	if err := os.Remove(absPath); err != nil {
		return NewInternalServerError(err)
	}

	return nil
}
