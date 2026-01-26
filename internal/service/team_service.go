package service

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

type TeamRepository interface {
	Create(ctx context.Context, t *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	List(ctx context.Context, limit, offset int, tenantID string) ([]domain.Team, error)
	Update(ctx context.Context, t *domain.Team) error
	Delete(ctx context.Context, id string) error
}

type TeamServicer interface {
	Create(ctx context.Context, t *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	List(ctx context.Context, limit, offset int, tenantID string) ([]domain.Team, error)
	Update(ctx context.Context, t *domain.Team) error
	Delete(ctx context.Context, id string) error
}

type TeamService struct {
	repo TeamRepository
}

func NewTeamService(repo TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) Create(ctx context.Context, t *domain.Team) error {
	return s.repo.Create(ctx, t)
}

func (s *TeamService) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TeamService) List(ctx context.Context, limit, offset int, tenantID string) ([]domain.Team, error) {
	return s.repo.List(ctx, limit, offset, tenantID)
}

func (s *TeamService) Update(ctx context.Context, t *domain.Team) error {
	return s.repo.Update(ctx, t)
}

func (s *TeamService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
