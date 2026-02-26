package middlewares

import (
	"fmt"
	"sync"
	"time"

	"phikhanh/utils"

	"github.com/gin-gonic/gin"
)

// uploadRecord - Lưu trữ thông tin upload của user
type uploadRecord struct {
	count     int
	totalSize int64
	resetAt   time.Time
}

var (
	uploadRecords = make(map[string]*uploadRecord)
	mu            sync.Mutex

	// Giới hạn upload trong 1 giờ
	maxUploadsPerHour   = 20
	maxTotalSizePerHour = int64(50 << 20) // 50MB per hour
	windowDuration      = time.Hour
)

// UploadRateLimitMiddleware - Middleware giới hạn số lần upload per user
func UploadRateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Dùng ExtractUserID
		userID, svcErr := utils.ExtractUserID(ctx)
		if svcErr != nil {
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			ctx.Abort()
			return
		}

		userIDStr := userID.String()

		// Lấy file size từ request
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.Next()
			return
		}

		mu.Lock()
		defer mu.Unlock()

		record, exists := uploadRecords[userIDStr]
		now := time.Now()

		// Reset record nếu đã qua thời gian window
		if !exists || now.After(record.resetAt) {
			uploadRecords[userIDStr] = &uploadRecord{
				count:     0,
				totalSize: 0,
				resetAt:   now.Add(windowDuration),
			}
			record = uploadRecords[userIDStr]
		}

		// Seconds còn lại cho đến khi reset
		retryAfter := int(time.Until(record.resetAt).Seconds())

		// Kiểm tra số lần upload
		if record.count >= maxUploadsPerHour {
			ctx.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			svcErr := utils.NewTooManyRequestsError(
				fmt.Sprintf("Upload limit exceeded (max %d uploads per hour). Retry after %d seconds", maxUploadsPerHour, retryAfter),
			)
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			ctx.Abort()
			return
		}

		// Kiểm tra tổng dung lượng
		if record.totalSize+file.Size > maxTotalSizePerHour {
			ctx.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			svcErr := utils.NewTooManyRequestsError(
				fmt.Sprintf("Upload size limit exceeded (max 50MB per hour). Retry after %d seconds", retryAfter),
			)
			utils.ErrorResponse(ctx, svcErr.StatusCode, svcErr.Message)
			ctx.Abort()
			return
		}

		// Cập nhật record
		record.count++
		record.totalSize += file.Size

		ctx.Next()
	}
}

// CleanupUploadRecords - Dọn dẹp records đã hết hạn (gọi định kỳ)
func CleanupUploadRecords() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	for userID, record := range uploadRecords {
		if now.After(record.resetAt) {
			delete(uploadRecords, userID)
		}
	}
}
