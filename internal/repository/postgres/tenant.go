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

type listTenantsRequest struct {
	ID        *string `json:"id"`
	Name      *string `json:"name"`
	LimitVal  int     `json:"limit_val"`
	OffsetVal int     `json:"offset_val"`
}

func (r *TenantRepo) List(ctx context.Context, limit, offset int) ([]domain.Tenant, error) {
	return ExecQueryList[domain.Tenant](ctx, r.pool, queryTenantList, listTenantsRequest{LimitVal: limit, OffsetVal: offset})
}
