package service

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

type TeamRepository interface {
	Create(ctx context.Context, t *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	List(ctx context.Context, limit, offset int, tenantID, userID string) ([]domain.Team, error)
	Update(ctx context.Context, t *domain.Team) error
	Delete(ctx context.Context, id string) error
	AddMember(ctx context.Context, teamID, userID, role string) (*domain.Member, error)
	IsMember(ctx context.Context, teamID, userID string) (bool, error)
}

type TeamServicer interface {
	Create(ctx context.Context, t *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	List(ctx context.Context, limit, offset int, tenantID, userID string) ([]domain.Team, error)
	Update(ctx context.Context, t *domain.Team) error
	Delete(ctx context.Context, id string) error
	AddMember(ctx context.Context, teamID, userID, role string) (*domain.Member, error)
	IsMember(ctx context.Context, teamID, userID string) (bool, error)
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

func (s *TeamService) List(ctx context.Context, limit, offset int, tenantID, userID string) ([]domain.Team, error) {
	return s.repo.List(ctx, limit, offset, tenantID, userID)
}

func (s *TeamService) Update(ctx context.Context, t *domain.Team) error {
	return s.repo.Update(ctx, t)
}

func (s *TeamService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *TeamService) AddMember(ctx context.Context, teamID, userID, role string) (*domain.Member, error) {
	return s.repo.AddMember(ctx, teamID, userID, role)
}

func (s *TeamService) IsMember(ctx context.Context, teamID, userID string) (bool, error) {
	return s.repo.IsMember(ctx, teamID, userID)
}
