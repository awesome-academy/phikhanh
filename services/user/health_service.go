package user

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}
