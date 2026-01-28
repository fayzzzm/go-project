package service

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

type TenantRepository interface {
	List(ctx context.Context, limit, offset int) ([]domain.Tenant, error)
}

type TenantService struct {
	repo TenantRepository
}

func NewTenantService(repo TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) List(ctx context.Context, limit, offset int) ([]domain.Tenant, error) {
	return s.repo.List(ctx, limit, offset)
}
