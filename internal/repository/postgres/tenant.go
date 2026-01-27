package postgres

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantRepo struct {
	pool *pgxpool.Pool
}

func NewTenantRepo(pool *pgxpool.Pool) *TenantRepo {
	return &TenantRepo{pool: pool}
}

const (
	queryTenantList = "SELECT id, name, created_at, updated_at FROM tenants.tenant ORDER BY created_at DESC LIMIT $1 OFFSET $2"
)

func (r *TenantRepo) List(ctx context.Context, limit, offset int) ([]domain.Tenant, error) {
	return ExecQueryList[domain.Tenant](ctx, r.pool, queryTenantList, limit, offset)
}
