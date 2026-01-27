package service

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

// UserRepository defined locally by the consumer.
type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context, limit, offset int, teamID, tenantID string) ([]domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id string) error
	GetForLogin(ctx context.Context, email string) (*domain.User, error)
}

type UserServicer interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context, limit, offset int, teamID, tenantID string) ([]domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, u *domain.User) error {
	return s.repo.Create(ctx, u)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, limit, offset int, teamID, tenantID string) ([]domain.User, error) {
	return s.repo.List(ctx, limit, offset, teamID, tenantID)
}

func (s *UserService) Update(ctx context.Context, u *domain.User) error {
	return s.repo.Update(ctx, u)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) GetForLogin(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetForLogin(ctx, email)
}
