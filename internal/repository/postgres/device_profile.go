package postgres

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceProfileRepo struct {
	pool *pgxpool.Pool
}

func NewDeviceProfileRepo(pool *pgxpool.Pool) *DeviceProfileRepo {
	return &DeviceProfileRepo{pool: pool}
}

const (
	queryDeviceProfileCreate  = "SELECT * FROM device_profiles.create($1::device_profiles.device_profile_request)"
	queryDeviceProfileGetByID = "SELECT * FROM device_profiles.get_by_id($1::device_profiles.device_profile_request)"
	queryDeviceProfileList    = "SELECT * FROM device_profiles.list($1::device_profiles.device_profile_request)"
	queryDeviceProfileUpdate  = "SELECT * FROM device_profiles.update($1::device_profiles.device_profile_request)"
	queryDeviceProfileDelete  = "SELECT * FROM device_profiles.delete($1::device_profiles.device_profile_request)"
)

func (r *DeviceProfileRepo) Create(ctx context.Context, dp *domain.DeviceProfile) error {
	val, err := ExecQueryOne[domain.DeviceProfile](ctx, r.pool, queryDeviceProfileCreate, NewDeviceProfileRequest(dp))
	if err == nil {
		*dp = *val
	}
	return err
}

func (r *DeviceProfileRepo) GetByID(ctx context.Context, id string) (*domain.DeviceProfile, error) {
	return ExecQueryOne[domain.DeviceProfile](ctx, r.pool, queryDeviceProfileGetByID, DeviceProfileRequest{ID: &id})
}

func (r *DeviceProfileRepo) Update(ctx context.Context, dp *domain.DeviceProfile) error {
	val, err := ExecQueryOne[domain.DeviceProfile](ctx, r.pool, queryDeviceProfileUpdate, NewDeviceProfileRequest(dp))
	if err == nil {
		*dp = *val
	}
	return err
}

func (r *DeviceProfileRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, queryDeviceProfileDelete, DeviceProfileRequest{ID: &id})
	return err
}

func (r *DeviceProfileRepo) List(ctx context.Context, limit, offset int) ([]domain.DeviceProfile, error) {
	return ExecQueryList[domain.DeviceProfile](ctx, r.pool, queryDeviceProfileList, DeviceProfileRequest{LimitVal: &limit, OffsetVal: &offset})
}
