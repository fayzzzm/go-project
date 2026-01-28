package postgres

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CabinetRepo struct {
	pool *pgxpool.Pool
}

func NewCabinetRepo(pool *pgxpool.Pool) *CabinetRepo {
	return &CabinetRepo{pool: pool}
}

const (
	queryCabinetCreate  = "SELECT * FROM cabinets.create($1::cabinets.cabinet_request)"
	queryCabinetGetByID = "SELECT * FROM cabinets.get_by_id($1::cabinets.cabinet_request)"
	queryCabinetList    = "SELECT * FROM cabinets.list($1::cabinets.cabinet_request)"
	queryCabinetUpdate  = "SELECT * FROM cabinets.update($1::cabinets.cabinet_request)"
	queryCabinetDelete  = "SELECT * FROM cabinets.delete($1::cabinets.cabinet_request)"
)

func (r *CabinetRepo) Create(ctx context.Context, c *domain.Cabinet, userID string) error {
	return ExecQueryUpdate(ctx, r.pool, c, queryCabinetCreate, NewCabinetRequest(c, &userID))
}

func (r *CabinetRepo) GetByID(ctx context.Context, id string) (*domain.Cabinet, error) {
	return ExecQueryOne[domain.Cabinet](ctx, r.pool, queryCabinetGetByID, &CabinetRequest{ID: &id})
}

func (r *CabinetRepo) Update(ctx context.Context, c *domain.Cabinet) error {
	return ExecQueryUpdate(ctx, r.pool, c, queryCabinetUpdate, NewCabinetRequest(c, nil))
}

func (r *CabinetRepo) Delete(ctx context.Context, id string) error {
	return Exec(ctx, r.pool, queryCabinetDelete, &CabinetRequest{ID: &id})
}

func (r *CabinetRepo) List(ctx context.Context, limit, offset int) ([]domain.Cabinet, error) {
	return ExecQueryList[domain.Cabinet](ctx, r.pool, queryCabinetList, &CabinetRequest{
		LimitVal: &limit, OffsetVal: &offset,
	})
}
