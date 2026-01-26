package service

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

// DeviceRepository is defined by the consumer (the service).
type DeviceRepository interface {
	Create(ctx context.Context, d *domain.Device) error
	GetByID(ctx context.Context, id string) (*domain.Device, error)
	List(ctx context.Context, limit, offset int) ([]domain.Device, error)
	GetBySerialNumber(ctx context.Context, serialNumber string) (*domain.Device, error)
	BulkCreate(ctx context.Context, devices []domain.Device) error
	Update(ctx context.Context, d *domain.Device) error
	Delete(ctx context.Context, id string) error
	GetByEPC(ctx context.Context, epc string) (*domain.Device, error)
}

type DeviceServicer interface {
	Create(ctx context.Context, d *domain.Device) error
	GetByID(ctx context.Context, id string) (*domain.Device, error)
	List(ctx context.Context, limit, offset int) ([]domain.Device, error)
	Update(ctx context.Context, d *domain.Device) error
	Delete(ctx context.Context, id string) error
	GetByEPC(ctx context.Context, epc string) (*domain.Device, error)
}

type DeviceService struct {
	repo DeviceRepository
}

func NewDeviceService(repo DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) Create(ctx context.Context, d *domain.Device) error {
	return s.repo.Create(ctx, d)
}

func (s *DeviceService) GetByID(ctx context.Context, id string) (*domain.Device, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DeviceService) List(ctx context.Context, limit, offset int) ([]domain.Device, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *DeviceService) Update(ctx context.Context, d *domain.Device) error {
	return s.repo.Update(ctx, d)
}

func (s *DeviceService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *DeviceService) GetByEPC(ctx context.Context, epc string) (*domain.Device, error) {
	return s.repo.GetByEPC(ctx, epc)
}
