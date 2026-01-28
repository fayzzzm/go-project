package service

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

type CabinetRepository interface {
	Create(ctx context.Context, c *domain.Cabinet, userID string) error
	GetByID(ctx context.Context, id string) (*domain.Cabinet, error)
	List(ctx context.Context, limit, offset int) ([]domain.Cabinet, error)
	Update(ctx context.Context, c *domain.Cabinet) error
	Delete(ctx context.Context, id string) error
}

type CabinetService struct {
	repo CabinetRepository
}

func NewCabinetService(repo CabinetRepository) *CabinetService {
	return &CabinetService{repo: repo}
}

func (s *CabinetService) Create(ctx context.Context, c *domain.Cabinet, userID string) error {
	return s.repo.Create(ctx, c, userID)
}

func (s *CabinetService) GetByID(ctx context.Context, id string) (*domain.Cabinet, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CabinetService) List(ctx context.Context, limit, offset int) ([]domain.Cabinet, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *CabinetService) Update(ctx context.Context, c *domain.Cabinet) error {
	return s.repo.Update(ctx, c)
}

func (s *CabinetService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
