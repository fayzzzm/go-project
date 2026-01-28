package usecase

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
)

type DeviceServicer interface {
	Create(ctx context.Context, d *domain.Device) error
	GetByID(ctx context.Context, id string) (*domain.Device, error)
	List(ctx context.Context, limit, offset int) ([]domain.Device, error)
	Update(ctx context.Context, d *domain.Device) error
	Delete(ctx context.Context, id string) error
	GetByEPC(ctx context.Context, epc string) (*domain.Device, error)
}

type CreateDeviceInput struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description"`
	SerialNumber    string `json:"serial_number" binding:"required"`
	EPC             string `json:"epc"`
	DeviceProfileID string `json:"device_profile_id"`
	CabinetID       string `json:"cabinet_id"`
	TeamID          string `json:"team_id"`
	TenantID        string `json:"tenant_id" binding:"required"`
}

type DeviceOutput struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
	Description  string `json:"description"`
	EPC          string `json:"epc"`
}

type DeviceUseCase struct {
	svc DeviceServicer
}

func NewDeviceUseCase(svc DeviceServicer) *DeviceUseCase {
	return &DeviceUseCase{svc: svc}
}

func (uc *DeviceUseCase) Create(ctx context.Context, input CreateDeviceInput) (*DeviceOutput, error) {
	device := &domain.Device{
		Name:            input.Name,
		Description:     utils.StringPtrOrNil(input.Description),
		SerialNumber:    input.SerialNumber,
		EPC:             utils.StringPtrOrNil(input.EPC),
		DeviceProfileID: utils.StringPtrOrNil(input.DeviceProfileID),
		CabinetID:       utils.StringPtrOrNil(input.CabinetID),
		TeamID:          utils.StringPtrOrNil(input.TeamID),
		TenantID:        utils.StringPtrOrNil(input.TenantID),
	}

	if err := uc.svc.Create(ctx, device); err != nil {
		return nil, err
	}

	return toDeviceOutput(device), nil
}

func (uc *DeviceUseCase) GetByID(ctx context.Context, id string) (*DeviceOutput, error) {
	device, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDeviceOutput(device), nil
}

func (uc *DeviceUseCase) List(ctx context.Context, p domain.Pagination) ([]DeviceOutput, error) {
	p.Normalize()
	devices, err := uc.svc.List(ctx, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}

	output := make([]DeviceOutput, len(devices))
	for i, d := range devices {
		output[i] = *toDeviceOutput(&d)
	}
	return output, nil
}

func (uc *DeviceUseCase) Update(ctx context.Context, id string, input CreateDeviceInput) (*DeviceOutput, error) {
	device := &domain.Device{
		ID:              id,
		Name:            input.Name,
		Description:     utils.StringPtrOrNil(input.Description),
		SerialNumber:    input.SerialNumber,
		EPC:             utils.StringPtrOrNil(input.EPC),
		DeviceProfileID: utils.StringPtrOrNil(input.DeviceProfileID),
		CabinetID:       utils.StringPtrOrNil(input.CabinetID),
		TeamID:          utils.StringPtrOrNil(input.TeamID),
		TenantID:        utils.StringPtrOrNil(input.TenantID),
	}

	if err := uc.svc.Update(ctx, device); err != nil {
		return nil, err
	}

	return toDeviceOutput(device), nil
}

func toDeviceOutput(d *domain.Device) *DeviceOutput {
	desc := ""
	if d.Description != nil {
		desc = *d.Description
	}
	epc := ""
	if d.EPC != nil {
		epc = *d.EPC
	}
	return &DeviceOutput{
		ID:           d.ID,
		Name:         d.Name,
		SerialNumber: d.SerialNumber,
		Description:  desc,
		EPC:          epc,
	}
}

func (uc *DeviceUseCase) Delete(ctx context.Context, id string) error {
	return uc.svc.Delete(ctx, id)
}
