package user

import (
	"fmt"
	"math"
	"time"

	userDto "phikhanh/dto/user"
	"phikhanh/models"
	userRepo "phikhanh/repositories/user"
	"phikhanh/utils"

	"github.com/google/uuid"
)

type ApplicationService struct {
	repo *userRepo.ApplicationRepository
}

func NewApplicationService(repo *userRepo.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
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
			Type:     models.AttachmentTypeRequest,
		})
	}

	// Tạo history record
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
