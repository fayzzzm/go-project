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
	return ExecQueryUpdate(ctx, r.pool, d, queryDeviceCreate, NewDeviceRequest(d))
}

func (r *DeviceRepo) GetByID(ctx context.Context, id string) (*domain.Device, error) {
	return ExecQueryOne[domain.Device](ctx, r.pool, queryDeviceGetByID, &DeviceRequest{ID: &id})
}

func (r *DeviceRepo) Update(ctx context.Context, d *domain.Device) error {
	return ExecQueryUpdate(ctx, r.pool, d, queryDeviceUpdate, NewDeviceRequest(d))
}

func (r *DeviceRepo) Delete(ctx context.Context, id string) error {
	return Exec(ctx, r.pool, queryDeviceDelete, &DeviceRequest{ID: &id})
}

func (r *DeviceRepo) List(ctx context.Context, limit, offset int) ([]domain.Device, error) {
	return ExecQueryList[domain.Device](ctx, r.pool, queryDeviceList, &DeviceRequest{
		LimitVal: &limit, OffsetVal: &offset,
	})
}

func (r *DeviceRepo) GetBySerialNumber(ctx context.Context, serialNumber string) (*domain.Device, error) {
	return ExecQueryOne[domain.Device](ctx, r.pool, queryDeviceGetBySerialNumber, &DeviceRequest{SerialNumber: &serialNumber})
}

func (r *DeviceRepo) GetByEPC(ctx context.Context, epc string) (*domain.Device, error) {
	return ExecQueryOne[domain.Device](ctx, r.pool, queryDeviceGetByEPC, &DeviceRequest{EPC: &epc})
}

func (r *DeviceRepo) BulkCreate(ctx context.Context, devices []domain.Device) error {
	requests := make([]*DeviceRequest, len(devices))
	for i, d := range devices {
		requests[i] = NewDeviceRequest(&d)
	}
	_, err := r.pool.Exec(ctx, queryDeviceBulkCreate, requests)
	return err
}
