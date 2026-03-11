package user

import (
	"fmt"
	"log"
	"math"
	"time"

	userDto "phikhanh/dto/user"
	"phikhanh/models"
	userRepo "phikhanh/repositories/user"
	"phikhanh/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationService struct {
	repo         *userRepo.ApplicationRepository
	emailService *utils.EmailService
}

func NewApplicationService(repo *userRepo.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		repo:         repo,
		emailService: utils.NewEmailService(),
	}
}

// SubmitApplication - Xử lý business logic nộp hồ sơ
func (s *ApplicationService) SubmitApplication(req userDto.SubmitAppRequest, userID uuid.UUID) (*userDto.SubmitAppResponse, error) {
	// Parse service_id
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, utils.NewBadRequestError("Invalid service_id format")
	}

	// Kiểm tra service có tồn tại không
	exists, err := s.repo.IsServiceExists(serviceID)
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}
	if !exists {
		return nil, utils.NewNotFoundError("Service not found")
	}

	// Generate unique application code
	code := generateApplicationCode()

	// Map DTO -> Application model
	app := &models.Application{
		Code:            code,
		UserID:          userID,
		ServiceID:       serviceID,
		Status:          models.StatusReceived,
		AssignedStaffID: nil, // Sẽ được assign sau khi Manager phân công
	}

	// Map DTO -> Attachments (optional, có thể rỗng)
	attachments := make([]models.Attachment, 0, len(req.Attachments))
	for _, a := range req.Attachments {
		attachments = append(attachments, models.Attachment{
			FileName: a.FileName,
			FilePath: a.FilePath,
			Type:     string(models.AttachmentTypeOriginal),
		})
	}

	// Tạo history record - ActorID là uuid.UUID (không phải pointer)
	history := &models.ApplicationHistory{
		ActorID: userID,
		Action:  "SUBMITTED",
		Note:    fmt.Sprintf("Application %s submitted by citizen", code),
	}

	// Lưu vào DB trong transaction
	if err := s.repo.CreateWithTransaction(app, attachments, history); err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	return &userDto.SubmitAppResponse{
		ID:   app.ID.String(),
		Code: app.Code,
	}, nil
}

// generateApplicationCode - Tạo mã hồ sơ unique
func generateApplicationCode() string {
	timestamp := time.Now().Format("20060102")
	shortID := uuid.New().String()[:8]
	return fmt.Sprintf("HS-%s-%s", timestamp, shortID)
}

// GetMyApplications - Lấy danh sách hồ sơ của user
func (s *ApplicationService) GetMyApplications(req userDto.MyAppListRequest, userID uuid.UUID) (*userDto.MyAppListResponse, error) {
	// Pagination defaults ĐÃ được set ở controller
	// Thêm defensive check để an toàn
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}

	applications, total, err := s.repo.FindMyApplications(userID, req)
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	items := make([]userDto.MyAppItemResponse, 0, len(applications))
	for _, app := range applications {
		item := userDto.MyAppItemResponse{
			ID:        app.ID.String(),
			Code:      app.Code,
			Status:    string(app.Status),
			CreatedAt: app.CreatedAt.Format(time.RFC3339),
		}

		if app.Service != nil {
			item.ServiceName = app.Service.Name
		}

		items = append(items, item)
	}

	// Safe: req.Limit > 0 vì đã validate ở controller
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &userDto.MyAppListResponse{
		Items:      items,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

// SupplementApplication - Validate + orchestrate supplement submission
func (s *ApplicationService) SupplementApplication(appID, userID string, req userDto.SupplementRequest) error {
	app, err := s.repo.FindByID(appID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewNotFoundError("Application not found")
		}
		return utils.NewInternalServerError(err)
	}

	if app.UserID.String() != userID {
		return utils.NewForbiddenError("Access denied: this application does not belong to you")
	}

	if app.Status != models.StatusSupplementRequired {
		return utils.NewBadRequestError("Application is not in a state that requires supplementation")
	}

	actorID, err := uuid.Parse(userID)
	if err != nil {
		return utils.NewBadRequestError("Invalid user ID")
	}

	appUUID, err := uuid.Parse(appID)
	if err != nil {
		return utils.NewBadRequestError("Invalid application ID")
	}

	attachments := make([]models.Attachment, 0, len(req.Attachments))
	for _, a := range req.Attachments {
		attachments = append(attachments, models.Attachment{
			ApplicationID: appUUID,
			FilePath:      a.FilePath,
			FileName:      a.FileName,
			Type:          string(models.AttachmentTypeSupplement),
		})
	}

	history := &models.ApplicationHistory{
		ApplicationID: appUUID,
		ActorID:       actorID,
		Action:        "SUPPLEMENTED",
		Note:          req.Note,
	}

	if err := s.repo.SupplementApplication(appID, userID, attachments, history); err != nil {
		return utils.NewInternalServerError(err)
	}

	// Notify assigned staff async
	go s.notifyStaffOnSupplementAsync(app, req.Note)

	return nil
}

// notifyStaffOnSupplementAsync - Gửi email thông báo cho staff khi citizen nộp bổ sung
func (s *ApplicationService) notifyStaffOnSupplementAsync(app *models.Application, citizenNote string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Supplement Email] Panic recovered: %v", r)
		}
	}()

	// Chỉ gửi nếu application có assigned staff
	if app.AssignedStaffID == nil {
		log.Printf("[Supplement Email] Skipping: no assigned staff for application %s", app.Code)
		return
	}

	// Fetch staff info
	staff, err := s.repo.FindUserByID(app.AssignedStaffID.String())
	if err != nil || staff == nil {
		log.Printf("[Supplement Email] Skipping: staff not found for application %s", app.Code)
		return
	}

	if staff.Email == "" {
		log.Printf("[Supplement Email] Skipping: staff email empty for application %s", app.Code)
		return
	}

	// Fetch citizen info
	citizenName := "Citizen"
	if app.User != nil {
		citizenName = app.User.Name
	}

	err = s.emailService.SendSupplementNotificationToStaff(
		staff.Email,
		staff.Name,
		app.Code,
		citizenName,
		citizenNote,
	)
	if err != nil {
		errorHandler := utils.GetEmailErrorHandler()
		errorHandler.ReportError(app.Code, staff.Email, err)
		return
	}

	log.Printf("[Supplement Email] ✓ Notified staff %s for application %s", staff.Email, app.Code)
}
