package admin

import (
	"fmt"
	"log"
	"math"
	adminDto "phikhanh/dto/admin"
	"phikhanh/models"
	adminRepo "phikhanh/repositories/admin"
	"phikhanh/services"
	"phikhanh/utils"
	"time"

	"gorm.io/gorm"
)

const (
	ApplicationPageSize = 10
)

type ApplicationAdminService struct {
	repo                *adminRepo.ApplicationRepository
	notificationService *services.NotificationService
	// emailService đã bị xóa: email được gửi qua notificationService.NotifyUser
}

func NewApplicationAdminService(repo *adminRepo.ApplicationRepository, notifService *services.NotificationService) *ApplicationAdminService {
	return &ApplicationAdminService{
		repo:                repo,
		notificationService: notifService,
	}
}

// GetList - Lấy danh sách applications với filter, pagination, và optional staff assignment filter
// Nếu assignedToUserID được cung cấp, chỉ lấy applications assigned to user này (dùng cho staff)
func (s *ApplicationAdminService) GetList(status string, assignedToUserID *string, page int) (*adminDto.ApplicationListResult, error) {
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * ApplicationPageSize

	applications, total, err := s.repo.FindAllWithFilterAndAssignment(status, assignedToUserID, offset, ApplicationPageSize)
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(ApplicationPageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	items := make([]adminDto.ApplicationListItem, 0, len(applications))
	for _, app := range applications {
		item := adminDto.ApplicationListItem{
			ID:          app.ID.String(),
			Code:        app.Code,
			Status:      string(app.Status),
			SubmittedAt: app.CreatedAt.Format(time.DateTime),
		}

		if app.User != nil {
			item.ApplicantName = app.User.Name
		}
		if app.Service != nil {
			item.ServiceName = app.Service.Name
		}

		items = append(items, item)
	}

	return &adminDto.ApplicationListResult{
		Items:       items,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		Status:      status,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}, nil
}

// GetDetail - Lấy chi tiết application với tất cả related data
func (s *ApplicationAdminService) GetDetail(id string) (*adminDto.ApplicationDetail, error) {
	app, err := s.repo.FindByIDWithDetailsAndDescription(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFoundError("Application not found")
		}
		return nil, utils.NewInternalServerError(err)
	}

	detail := &adminDto.ApplicationDetail{
		ID:          app.ID.String(),
		Code:        app.Code,
		Status:      string(app.Status),
		SubmittedAt: app.CreatedAt.Format(time.DateTime),
	}

	// Map AssignedStaffID để check authorization
	if app.AssignedStaffID != nil {
		detail.AssignedStaffID = app.AssignedStaffID.String()
	}

	if app.User != nil {
		detail.ApplicantName = app.User.Name
		detail.CitizenID = app.User.CitizenID
		detail.Email = app.User.Email
		detail.Phone = app.User.Phone
	}

	if app.Service != nil {
		detail.ServiceName = app.Service.Name
		detail.ProcessingDays = app.Service.ProcessingDays
		detail.Fee = app.Service.Fee
	}

	// Map attachments
	detail.Attachments = make([]adminDto.ApplicationAttachment, 0, len(app.Attachments))
	for _, att := range app.Attachments {
		detail.Attachments = append(detail.Attachments, adminDto.ApplicationAttachment{
			FileName: att.FileName,
			FilePath: att.FilePath,
		})
	}

	// Map histories
	detail.Histories = make([]adminDto.ApplicationHistory, 0, len(app.Histories))
	for _, hist := range app.Histories {
		h := adminDto.ApplicationHistory{
			Date:        hist.CreatedAt.Format(time.DateTime),
			Action:      hist.Action,
			Note:        hist.Note,
			Description: hist.Note,
		}
		if hist.Actor != nil {
			h.ActorName = hist.Actor.Name
		}
		detail.Histories = append(detail.Histories, h)
	}

	return detail, nil
}

// GetAvailableStaff - Lấy danh sách staff (role = "staff") để assign applications
func (s *ApplicationAdminService) GetAvailableStaff() ([]adminDto.StaffMember, error) {
	staff, err := s.repo.GetAvailableStaff()
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	result := make([]adminDto.StaffMember, 0, len(staff))
	for _, u := range staff {
		result = append(result, adminDto.StaffMember{
			ID:   u.ID.String(),
			Name: u.Name,
			Role: string(u.Role),
		})
	}

	return result, nil
}

// GetNextStatuses - Lấy danh sách status tiếp theo dựa trên status hiện tại
// Mỗi status có workflow state machine với các transition cho phép
func (s *ApplicationAdminService) GetNextStatuses(currentStatus string) []string {
	// Status Workflow State Machine
	// Received -> Processing, Supplement_Required, Rejected
	// Processing -> Approved, Supplement_Required, Rejected
	// Supplement_Required -> Processing, Approved, Rejected
	// Approved -> (terminal state)
	// Rejected -> (terminal state)
	statusTransitions := map[string][]string{
		string(models.StatusReceived): {
			string(models.StatusProcessing),
			string(models.StatusSupplementRequired),
			string(models.StatusRejected),
		},
		string(models.StatusProcessing): {
			string(models.StatusApproved),
			string(models.StatusSupplementRequired),
			string(models.StatusRejected),
		},
		string(models.StatusSupplementRequired): {
			string(models.StatusProcessing),
			string(models.StatusApproved),
			string(models.StatusRejected),
		},
		string(models.StatusApproved): {},
		string(models.StatusRejected): {},
	}

	if statuses, ok := statusTransitions[currentStatus]; ok {
		return statuses
	}

	return []string{}
}

// ProcessApplication - Main workflow để xử lý application
// 1. Validate status transition dựa trên state machine
// 2. Get current app data
// 3. Build history description
// 4. Update status + assign staff (transaction)
// 5. Trigger async email notification
func (s *ApplicationAdminService) ProcessApplication(appID string, newStatus string, assignedStaffID *string, note string, actorID string) error {
	// Lấy current application để biết old status
	app, err := s.repo.FindByIDWithDetailsAndDescription(appID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewNotFoundError("Application not found")
		}
		return utils.NewInternalServerError(err)
	}

	oldStatus := string(app.Status)

	// Validate status transition dựa trên state machine
	allowedTransitions := s.GetNextStatuses(oldStatus)
	if !s.isAllowedTransition(newStatus, allowedTransitions) {
		return utils.NewBadRequestError(
			"Invalid status transition: cannot change from " + oldStatus + " to " + newStatus,
		)
	}

	// Lấy tên staff nếu được assign
	var assignedStaffName *string
	if assignedStaffID != nil && *assignedStaffID != "" {
		name, err := s.repo.GetStaffNameByID(*assignedStaffID)
		if err == nil && name != "" {
			assignedStaffName = &name
		}
	}

	// Build description chi tiết cho history
	description := s.BuildHistoryDescription(oldStatus, newStatus, assignedStaffName, note)

	// Update application trong DB + insert history record (transaction)
	if err := s.repo.ProcessAndAssignWithHistoryV2(appID, oldStatus, newStatus, assignedStaffID, assignedStaffName, note, actorID, description); err != nil {
		return utils.NewInternalServerError(err)
	}

	// Notify citizen async:
	// - Luôn tạo in-app notification
	// - Gửi email nếu user.IsEmailNotify = true (handled bên trong NotifyUser)
	if app.User != nil {
		go s.notifyCitizenOnStatusChange(app, newStatus, note)
	}

	return nil
}

// isAllowedTransition - Check nếu transition từ oldStatus tới newStatus có được phép không
func (s *ApplicationAdminService) isAllowedTransition(newStatus string, allowedStatuses []string) bool {
	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

// sendStatusUpdateEmailAsync - DEPRECATED: đã được thay bởi notifyCitizenOnStatusChange + NotifyUser
// Giữ lại để tham khảo, không gọi nữa
// func (s *ApplicationAdminService) sendStatusUpdateEmailAsync(...) { ... }

// BuildHistoryDescription - Xây dựng mô tả chi tiết cho history record
// Format: "Status changed from Received to Processing and assigned to Nguyen Van A. Note: ..."
func (s *ApplicationAdminService) BuildHistoryDescription(oldStatus string, newStatus string, assignedStaffName *string, note string) string {
	desc := fmt.Sprintf("Status changed from %s to %s", oldStatus, newStatus)

	if assignedStaffName != nil && *assignedStaffName != "" {
		desc += fmt.Sprintf(" and assigned to %s", *assignedStaffName)
	}

	if note != "" {
		desc += fmt.Sprintf(". Note: %s", note)
	}

	return desc
}

// notifyCitizenOnStatusChange - Delegate toàn bộ notification + email logic cho NotificationService
// NotificationService.NotifyUser sẽ:
//   - Step 1: Insert in-app notification
//   - Step 2: Fetch user.IsEmailNotify
//   - Step 3: Gửi email async nếu IsEmailNotify = true
func (s *ApplicationAdminService) notifyCitizenOnStatusChange(app *models.Application, newStatus, note string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Notification] Panic recovered: %v", r)
		}
	}()

	title := fmt.Sprintf("Hồ sơ %s - Trạng thái: %s", app.Code, newStatus)
	content := fmt.Sprintf("Hồ sơ của bạn đã được cập nhật sang trạng thái: %s", newStatus)
	if note != "" {
		content += fmt.Sprintf("\nGhi chú từ cán bộ: %s", note)
	}

	if err := s.notificationService.NotifyUser(app.UserID, title, content); err != nil {
		log.Printf("[Notification] Failed to notify citizen %s for app %s: %v", app.UserID, app.Code, err)
	}
}
