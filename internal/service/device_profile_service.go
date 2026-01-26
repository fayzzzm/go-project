package service

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

// DeviceProfileRepository defined locally by consumer
type DeviceProfileRepository interface {
	Create(ctx context.Context, dp *domain.DeviceProfile) error
	GetByID(ctx context.Context, id string) (*domain.DeviceProfile, error)
	List(ctx context.Context, limit, offset int) ([]domain.DeviceProfile, error)
	Update(ctx context.Context, dp *domain.DeviceProfile) error
	Delete(ctx context.Context, id string) error
}

type DeviceProfileServicer interface {
	Create(ctx context.Context, dp *domain.DeviceProfile) error
	GetByID(ctx context.Context, id string) (*domain.DeviceProfile, error)
	List(ctx context.Context, limit, offset int) ([]domain.DeviceProfile, error)
	Update(ctx context.Context, dp *domain.DeviceProfile) error
	Delete(ctx context.Context, id string) error
}

type DeviceProfileService struct {
	repo DeviceProfileRepository
}

func NewDeviceProfileService(repo DeviceProfileRepository) *DeviceProfileService {
	return &DeviceProfileService{repo: repo}
}

func (s *DeviceProfileService) Create(ctx context.Context, dp *domain.DeviceProfile) error {
	return s.repo.Create(ctx, dp)
}

func (s *DeviceProfileService) GetByID(ctx context.Context, id string) (*domain.DeviceProfile, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DeviceProfileService) List(ctx context.Context, limit, offset int) ([]domain.DeviceProfile, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *DeviceProfileService) Update(ctx context.Context, dp *domain.DeviceProfile) error {
	return s.repo.Update(ctx, dp)
}

func (s *DeviceProfileService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
