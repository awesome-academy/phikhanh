package services

import (
	"log"
	"math"

	"phikhanh/models"
	"phikhanh/repositories"
	"phikhanh/utils"

	"github.com/google/uuid"
)

const NotificationPageSize = 20

type NotificationService struct {
	repo         *repositories.NotificationRepository
	emailService *utils.EmailService
}

func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo:         repo,
		emailService: utils.NewEmailService(),
	}
}

// NotifyUser - Luôn tạo in-app notification.
// Sau đó check IsEmailNotify: nếu true → gửi email async.
func (s *NotificationService) NotifyUser(userID uuid.UUID, title, content string) error {
	// Step 1: Insert in-app notification (luôn thực hiện)
	notification := &models.Notification{
		UserID:  userID,
		Title:   title,
		Content: content,
		IsRead:  false,
	}
	if err := s.repo.Create(notification); err != nil {
		return utils.NewInternalServerError(err)
	}

	// Step 2 & 3: Check IsEmailNotify và gửi email async (không block)
	go s.sendEmailIfEnabled(userID, title, content)

	return nil
}

// sendEmailIfEnabled - Fetch user preference và gửi email nếu IsEmailNotify = true
func (s *NotificationService) sendEmailIfEnabled(userID uuid.UUID, title, content string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Notification Email] Panic recovered for user %s: %v", userID, r)
		}
	}()

	// Step 2: Fetch user với chỉ fields cần thiết (email + IsEmailNotify)
	user, err := s.repo.FindUserWithEmailPref(userID)
	if err != nil {
		log.Printf("[Notification Email] Failed to fetch user %s: %v", userID, err)
		return
	}

	// Step 3: Chỉ gửi email nếu user bật IsEmailNotify
	if !user.IsEmailNotify {
		log.Printf("[Notification Email] Skipping: user %s has IsEmailNotify=false", userID)
		return
	}

	if user.Email == "" {
		log.Printf("[Notification Email] Skipping: user %s has empty email", userID)
		return
	}

	if !s.emailService.IsConfigured() {
		log.Printf("[Notification Email] Skipping: SMTP not configured")
		return
	}

	// Gửi email thông báo - dùng title làm subject, content làm body note
	err = s.emailService.SendNotificationEmail(user.Email, user.Name, title, content)
	if err != nil {
		errorHandler := utils.GetEmailErrorHandler()
		errorHandler.ReportError(userID.String(), user.Email, err)
		return
	}

	log.Printf("[Notification Email] ✓ Email sent to %s for notification: %s", user.Email, title)
}

// GetNotifications - Lấy danh sách notifications của user với pagination
func (s *NotificationService) GetNotifications(userID uuid.UUID, page int) (*NotificationListResult, error) {
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * NotificationPageSize

	notifications, total, err := s.repo.FindByUserID(userID, offset, NotificationPageSize)
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(NotificationPageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &NotificationListResult{
		Items:       notifications,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}, nil
}

// MarkAsRead - Đánh dấu 1 notification đã đọc
func (s *NotificationService) MarkAsRead(notificationID, userID uuid.UUID) error {
	if err := s.repo.MarkAsRead(notificationID, userID); err != nil {
		return utils.NewNotFoundError("Notification not found")
	}
	return nil
}

// MarkAllAsRead - Đánh dấu tất cả notifications của user đã đọc
func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.repo.MarkAllAsRead(userID)
}

// CountUnread - Đếm số notifications chưa đọc
func (s *NotificationService) CountUnread(userID uuid.UUID) (int64, error) {
	return s.repo.CountUnread(userID)
}

// NotificationListResult - Pagination result
type NotificationListResult struct {
	Items       []models.Notification `json:"items"`
	CurrentPage int                   `json:"current_page"`
	TotalPages  int                   `json:"total_pages"`
	TotalItems  int64                 `json:"total_items"`
	HasPrev     bool                  `json:"has_prev"`
	HasNext     bool                  `json:"has_next"`
}
