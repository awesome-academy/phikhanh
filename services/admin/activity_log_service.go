package admin

import (
	"log"
	"math"
	"time"

	adminDto "phikhanh/dto/admin"
	"phikhanh/models"
	adminRepo "phikhanh/repositories/admin"
	"phikhanh/utils"

	"github.com/google/uuid"
)

const ActivityLogPageSize = 20

type ActivityLogService struct {
	repo *adminRepo.ActivityLogRepository
}

func NewActivityLogService(repo *adminRepo.ActivityLogRepository) *ActivityLogService {
	return &ActivityLogService{repo: repo}
}

// RecordActivity - Reusable helper, gọi async từ bất kỳ controller nào
func (s *ActivityLogService) RecordActivity(actorID, action, targetID, description, ipAddress string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ActivityLog] Panic: %v", r)
			}
		}()

		logEntry := &models.SystemLog{
			Action:      action,
			TargetID:    targetID,
			Description: description,
			IPAddress:   ipAddress,
		}

		if actorID != "" {
			if id, err := uuid.Parse(actorID); err == nil {
				logEntry.ActorID = &id
			}
		}

		if err := s.repo.Create(logEntry); err != nil {
			log.Printf("[ActivityLog] Failed to record: action=%s, err=%v", action, err)
		}
	}()
}

func (s *ActivityLogService) GetList(action, keyword string, page int) (*adminDto.ActivityLogListResult, error) {
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * ActivityLogPageSize

	logs, total, err := s.repo.FindAllWithFilter(action, keyword, offset, ActivityLogPageSize)
	if err != nil {
		return nil, utils.NewInternalServerError(err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(ActivityLogPageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	items := make([]adminDto.ActivityLogListItem, 0, len(logs))
	for _, l := range logs {
		item := adminDto.ActivityLogListItem{
			ID:          l.ID.String(),
			Action:      l.Action,
			TargetID:    l.TargetID,
			Description: l.Description,
			IPAddress:   l.IPAddress,
			CreatedAt:   l.CreatedAt.Format(time.DateTime),
			ActorName:   "System",
		}
		if l.Actor != nil {
			item.ActorName = l.Actor.Name
		}
		items = append(items, item)
	}

	return &adminDto.ActivityLogListResult{
		Items:       items,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		Action:      action,
		Keyword:     keyword,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}, nil
}

func (s *ActivityLogService) CleanupOldLogs(days int) (int64, error) {
	if days <= 0 {
		days = 30
	}
	count, err := s.repo.DeleteOlderThan(days)
	if err != nil {
		return 0, utils.NewInternalServerError(err)
	}
	return count, nil
}

func (s *ActivityLogService) GetAvailableActions() []string {
	return s.repo.GetAvailableActions()
}
