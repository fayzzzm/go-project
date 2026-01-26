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

func (r *CabinetRepo) Create(ctx context.Context, c *domain.Cabinet) error {
	// Status is string in domain, *string in Request.
	var status *string
	if c.Status != "" {
		val := c.Status
		status = &val
	}

	req := CabinetRequest{
		Name:        &c.Name,
		Description: c.Description, // Already *string
		Location:    c.Location,    // Already *string
		MachineID:   c.MachineID,   // Already *string
		Status:      status,
		TeamID:      c.TeamID,   // Already *string
		TenantID:    c.TenantID, // Already *string
	}
	val, err := ExecQueryOne[domain.Cabinet](ctx, r.pool, queryCabinetCreate, req)
	if err != nil {
		return err
	}
	*c = *val
	return nil
}

func (r *CabinetRepo) GetByID(ctx context.Context, id string) (*domain.Cabinet, error) {
	req := CabinetRequest{ID: &id}
	return ExecQueryOne[domain.Cabinet](ctx, r.pool, queryCabinetGetByID, req)
}

func (r *CabinetRepo) Update(ctx context.Context, c *domain.Cabinet) error {
	// Status is string in domain, *string in Request.
	var status *string
	if c.Status != "" {
		val := c.Status
		status = &val
	}

	req := CabinetRequest{
		ID:          &c.ID,
		Name:        &c.Name,
		Description: c.Description, // Already *string
		Location:    c.Location,    // Already *string
		MachineID:   c.MachineID,   // Already *string
		Status:      status,
		TeamID:      c.TeamID,
	}
	val, err := ExecQueryOne[domain.Cabinet](ctx, r.pool, queryCabinetUpdate, req)
	if err != nil {
		return err
	}
	*c = *val
	return nil
}

func (r *CabinetRepo) Delete(ctx context.Context, id string) error {
	req := CabinetRequest{ID: &id}
	_, err := r.pool.Exec(ctx, queryCabinetDelete, req)
	return err
}

func (r *CabinetRepo) List(ctx context.Context, limit, offset int) ([]domain.Cabinet, error) {
	req := CabinetRequest{LimitVal: &limit, OffsetVal: &offset}
	return ExecQueryList[domain.Cabinet](ctx, r.pool, queryCabinetList, req)
}
