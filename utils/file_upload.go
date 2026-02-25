package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	UploadDir   = "./assets/images"
	MaxFileSize = 10 << 20 // 10MB
	AllowedExts = ".jpg,.jpeg,.png,.pdf,.doc,.docx"
)

type UploadResult struct {
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	FileURL      string `json:"file_url"`
	FileSize     int64  `json:"file_size"`
	OriginalName string `json:"original_name"`
}

// UploadFile - Upload file và trả về đường dẫn
func UploadFile(file *multipart.FileHeader) (*UploadResult, error) {
	// Validate file size
	if file.Size > MaxFileSize {
		return nil, NewBadRequestError("File size exceeds maximum allowed (10MB)")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isAllowedExtension(ext) {
		return nil, NewBadRequestError(fmt.Sprintf("File type not allowed. Allowed types: %s", AllowedExts))
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		return nil, NewInternalServerError(err)
	}

	// Generate unique filename
	uniqueFileName := generateUniqueFileName(file.Filename)
	filePath := filepath.Join(UploadDir, uniqueFileName)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, NewInternalServerError(err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, NewInternalServerError(err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := dst.ReadFrom(src); err != nil {
		return nil, NewInternalServerError(err)
	}

	// Generate file URL (relative path for API response)
	fileURL := fmt.Sprintf("/assets/images/%s", uniqueFileName)

	return &UploadResult{
		FileName:     uniqueFileName,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     file.Size,
		OriginalName: file.Filename,
	}, nil
}

// generateUniqueFileName - Tạo tên file unique
func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
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

// DeleteFile - Xóa file từ server
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return NewNotFoundError("File not found")
	}

	if err := os.Remove(filePath); err != nil {
		return NewInternalServerError(err)
	}

	return nil
}
