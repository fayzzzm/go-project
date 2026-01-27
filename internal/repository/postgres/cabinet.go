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
	val, err := ExecQueryOne[domain.Cabinet](ctx, r.pool, queryCabinetCreate, NewCabinetRequest(c, &userID))
	if err == nil {
		*c = *val
	}
	return err
}

func (r *CabinetRepo) GetByID(ctx context.Context, id string) (*domain.Cabinet, error) {
	return ExecQueryOne[domain.Cabinet](ctx, r.pool, queryCabinetGetByID, CabinetRequest{ID: &id})
}

func (r *CabinetRepo) Update(ctx context.Context, c *domain.Cabinet) error {
	val, err := ExecQueryOne[domain.Cabinet](ctx, r.pool, queryCabinetUpdate, NewCabinetRequest(c, nil))
	if err == nil {
		*c = *val
	}
	return err
}

func (r *CabinetRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, queryCabinetDelete, CabinetRequest{ID: &id})
	return err
}

func (r *CabinetRepo) List(ctx context.Context, limit, offset int) ([]domain.Cabinet, error) {
	return ExecQueryList[domain.Cabinet](ctx, r.pool, queryCabinetList, CabinetRequest{LimitVal: &limit, OffsetVal: &offset})
}
