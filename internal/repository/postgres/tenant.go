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
	queryTenantList = "SELECT * FROM tenants.list($1)"
)

func (r *TenantRepo) List(ctx context.Context, limit, offset int) ([]domain.Tenant, error) {
	return ExecQueryList[domain.Tenant](ctx, r.pool, queryTenantList, &TenantRequest{
		LimitVal: &limit, OffsetVal: &offset,
	})
}
