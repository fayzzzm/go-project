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
	Name            string  `json:"name" binding:"required"`
	Description     *string `json:"description"`
	SerialNumber    string  `json:"serial_number" binding:"required"`
	EPC             *string `json:"epc"`
	DeviceProfileID *string `json:"device_profile_id"`
	CabinetID       *string `json:"cabinet_id"`
	TeamID          *string `json:"team_id"`
}

type UpdateDeviceInput struct {
	ID              string  `json:"-"`
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	SerialNumber    *string `json:"serial_number"`
	EPC             *string `json:"epc"`
	DeviceProfileID *string `json:"device_profile_id"`
	CabinetID       *string `json:"cabinet_id"`
	TeamID          *string `json:"team_id"`
}

type DeviceIDInput struct {
	ID string `json:"-"`
}

type ListDeviceInput struct {
	Pagination domain.Pagination
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
	userID := domain.UserFromContext(ctx)

	device := &domain.Device{
		Name:            input.Name,
		Description:     input.Description,
		SerialNumber:    input.SerialNumber,
		EPC:             input.EPC,
		DeviceProfileID: input.DeviceProfileID,
		CabinetID:       input.CabinetID,
		TeamID:          input.TeamID,
		CreatedBy:       utils.StringPtrOrNil(userID),
		UpdatedBy:       utils.StringPtrOrNil(userID),
	}

	if err := uc.svc.Create(ctx, device); err != nil {
		return nil, err
	}

	return toDeviceOutput(device), nil
}

func (uc *DeviceUseCase) GetByID(ctx context.Context, input DeviceIDInput) (*DeviceOutput, error) {
	device, err := uc.svc.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return toDeviceOutput(device), nil
}

func (uc *DeviceUseCase) List(ctx context.Context, input ListDeviceInput) ([]DeviceOutput, error) {
	input.Pagination.Normalize()
	devices, err := uc.svc.List(ctx, input.Pagination.Limit, input.Pagination.Offset)
	if err != nil {
		return nil, err
	}

	output := make([]DeviceOutput, len(devices))
	for i, d := range devices {
		output[i] = *toDeviceOutput(&d)
	}
	return output, nil
}

func (uc *DeviceUseCase) Update(ctx context.Context, input UpdateDeviceInput) (*DeviceOutput, error) {
	device, err := uc.svc.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		device.Name = *input.Name
	}
	if input.Description != nil {
		device.Description = input.Description
	}
	if input.SerialNumber != nil {
		device.SerialNumber = *input.SerialNumber
	}
	if input.EPC != nil {
		device.EPC = input.EPC
	}
	if input.DeviceProfileID != nil {
		device.DeviceProfileID = input.DeviceProfileID
	}
	if input.CabinetID != nil {
		device.CabinetID = input.CabinetID
	}
	if input.TeamID != nil {
		device.TeamID = input.TeamID
	}

	if userID := domain.UserFromContext(ctx); userID != "" {
		device.UpdatedBy = utils.StringPtrOrNil(userID)
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

func (uc *DeviceUseCase) Delete(ctx context.Context, input DeviceIDInput) error {
	return uc.svc.Delete(ctx, input.ID)
}
