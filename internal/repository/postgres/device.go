package postgres

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepo struct {
	pool *pgxpool.Pool
}

func NewDeviceRepo(pool *pgxpool.Pool) *DeviceRepo {
	return &DeviceRepo{pool: pool}
}

const (
	queryDeviceCreate            = "SELECT * FROM devices.create($1)"
	queryDeviceGetByID           = "SELECT * FROM devices.get_by_id($1)"
	queryDeviceList              = "SELECT * FROM devices.list($1)"
	queryDeviceUpdate            = "SELECT * FROM devices.update($1)"
	queryDeviceDelete            = "SELECT devices.delete($1)"
	queryDeviceGetBySerialNumber = "SELECT * FROM devices.get_by_serial_number($1)"
	queryDeviceBulkCreate        = "SELECT * FROM devices.bulk_create($1)"
	queryDeviceGetByEPC          = "SELECT * FROM devices.get_by_epc($1)"
)

func (r *DeviceRepo) Create(ctx context.Context, d *domain.Device) error {
	req := DeviceRequest{
		Name:            &d.Name,
		Description:     &d.Description,
		SerialNumber:    &d.SerialNumber,
		EPC:             &d.EPC,
		DeviceProfileID: d.DeviceProfileID,
		CabinetID:       d.CabinetID,
		TeamID:          d.TeamID,
		TenantID:        &d.TenantID,
	}
	val, err := ExecQueryOne[domain.Device](ctx, r.pool, queryDeviceCreate, req)
	if err != nil {
		return err
	}
	*d = *val
	return nil
}

func (r *DeviceRepo) GetByID(ctx context.Context, id string) (*domain.Device, error) {
	req := DeviceRequest{ID: &id}
	return ExecQueryOne[domain.Device](ctx, r.pool, queryDeviceGetByID, req)
}

func (r *DeviceRepo) Update(ctx context.Context, d *domain.Device) error {
	req := DeviceRequest{
		ID:              &d.ID,
		Name:            &d.Name,
		Description:     &d.Description,
		SerialNumber:    &d.SerialNumber,
		EPC:             &d.EPC,
		DeviceProfileID: d.DeviceProfileID,
		CabinetID:       d.CabinetID,
		TeamID:          d.TeamID,
		TenantID:        &d.TenantID,
	}
	val, err := ExecQueryOne[domain.Device](ctx, r.pool, queryDeviceUpdate, req)
	if err != nil {
		return err
	}
	*d = *val
	return nil
}

func (r *DeviceRepo) Delete(ctx context.Context, id string) error {
	req := DeviceRequest{ID: &id}
	_, err := r.pool.Exec(ctx, queryDeviceDelete, req)
	return err
}

func (r *DeviceRepo) List(ctx context.Context, limit, offset int) ([]domain.Device, error) {
	req := DeviceRequest{LimitVal: &limit, OffsetVal: &offset}
	return ExecQueryList[domain.Device](ctx, r.pool, queryDeviceList, req)
}

func (r *DeviceRepo) GetBySerialNumber(ctx context.Context, serialNumber string) (*domain.Device, error) {
	req := DeviceRequest{SerialNumber: &serialNumber}
	return ExecQueryOne[domain.Device](ctx, r.pool, queryDeviceGetBySerialNumber, req)
}

func (r *DeviceRepo) GetByEPC(ctx context.Context, epc string) (*domain.Device, error) {
	req := DeviceRequest{EPC: &epc}
	return ExecQueryOne[domain.Device](ctx, r.pool, queryDeviceGetByEPC, req)
}

func (r *DeviceRepo) BulkCreate(ctx context.Context, devices []domain.Device) error {
	requests := make([]DeviceRequest, len(devices))
	for i, d := range devices {
		requests[i] = DeviceRequest{
			Name:            &d.Name,
			Description:     &d.Description,
			SerialNumber:    &d.SerialNumber,
			EPC:             &d.EPC,
			DeviceProfileID: d.DeviceProfileID,
			CabinetID:       d.CabinetID,
			TeamID:          d.TeamID,
		}
	}
	_, err := r.pool.Exec(ctx, queryDeviceBulkCreate, requests)
	return err
}
