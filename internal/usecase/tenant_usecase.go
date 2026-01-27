package usecase

import (
	"context"
	"time"

	"github.com/fayzzzm/go-project/internal/domain"
)

type TenantServicer interface {
	List(ctx context.Context, limit, offset int) ([]domain.Tenant, error)
}

type TenantOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantUseCase struct {
	svc TenantServicer
}

func NewTenantUseCase(svc TenantServicer) *TenantUseCase {
	return &TenantUseCase{svc: svc}
}

func (uc *TenantUseCase) List(ctx context.Context, p domain.Pagination) ([]TenantOutput, error) {
	p.Normalize()
	tenants, err := uc.svc.List(ctx, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}

	output := make([]TenantOutput, len(tenants))
	for i, t := range tenants {
		output[i] = TenantOutput{
			ID:        t.ID,
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		}
	}
	return output, nil
}
