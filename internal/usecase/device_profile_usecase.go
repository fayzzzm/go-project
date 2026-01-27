package usecase

import (
	"context"
	"time"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
)

type DeviceProfileServicer interface {
	Create(ctx context.Context, dp *domain.DeviceProfile) error
	GetByID(ctx context.Context, id string) (*domain.DeviceProfile, error)
	List(ctx context.Context, limit, offset int) ([]domain.DeviceProfile, error)
	Update(ctx context.Context, dp *domain.DeviceProfile) error
	Delete(ctx context.Context, id string) error
}

type CreateDeviceProfileInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	TenantID    string `json:"tenant_id"`
}

type UpdateDeviceProfileInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeviceProfileOutput struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TenantID    string    `json:"tenant_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeviceProfileUseCase struct {
	svc DeviceProfileServicer
}

func NewDeviceProfileUseCase(svc DeviceProfileServicer) *DeviceProfileUseCase {
	return &DeviceProfileUseCase{svc: svc}
}

func (uc *DeviceProfileUseCase) Create(ctx context.Context, input CreateDeviceProfileInput) (*DeviceProfileOutput, error) {
	dp := &domain.DeviceProfile{
		Name:        input.Name,
		Description: utils.StringPtrOrNil(input.Description),
		TenantID:    utils.StringPtrOrNil(input.TenantID),
	}

	if err := uc.svc.Create(ctx, dp); err != nil {
		return nil, err
	}

	return toDeviceProfileOutput(dp), nil
}

func (uc *DeviceProfileUseCase) GetByID(ctx context.Context, id string) (*DeviceProfileOutput, error) {
	dp, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDeviceProfileOutput(dp), nil
}

func (uc *DeviceProfileUseCase) List(ctx context.Context, p domain.Pagination) ([]DeviceProfileOutput, error) {
	p.Normalize()
	dps, err := uc.svc.List(ctx, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}

	output := make([]DeviceProfileOutput, len(dps))
	for i, dp := range dps {
		output[i] = *toDeviceProfileOutput(&dp)
	}
	return output, nil
}

func (uc *DeviceProfileUseCase) Update(ctx context.Context, id string, input UpdateDeviceProfileInput) (*DeviceProfileOutput, error) {
	dp := &domain.DeviceProfile{
		ID:          id,
		Name:        input.Name,
		Description: utils.StringPtrOrNil(input.Description),
	}

	if err := uc.svc.Update(ctx, dp); err != nil {
		return nil, err
	}
	return toDeviceProfileOutput(dp), nil
}

func (uc *DeviceProfileUseCase) Delete(ctx context.Context, id string) error {
	return uc.svc.Delete(ctx, id)
}

func toDeviceProfileOutput(dp *domain.DeviceProfile) *DeviceProfileOutput {
	return &DeviceProfileOutput{
		ID:          dp.ID,
		Name:        dp.Name,
		Description: utils.StringValue(dp.Description),
		TenantID:    utils.StringValue(dp.TenantID),
		CreatedAt:   dp.CreatedAt,
		UpdatedAt:   dp.UpdatedAt,
	}
}
