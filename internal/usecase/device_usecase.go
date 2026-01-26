package usecase

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
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
	Name            string `json:"name"`
	Description     string `json:"description"`
	SerialNumber    string `json:"serial_number"`
	EPC             string `json:"epc"`
	DeviceProfileID string `json:"device_profile_id"`
	CabinetID       string `json:"cabinet_id"`
	TeamID          string `json:"team_id"`
	TenantID        string `json:"tenant_id"`
}

func (i CreateDeviceInput) Validate() (bool, string) {
	d := domain.Device{
		Name:         i.Name,
		SerialNumber: i.SerialNumber,
		TenantID:     i.TenantID,
	}
	return d.Validate()
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
	var profileID, cabinetID, teamID *string
	if input.DeviceProfileID != "" {
		val := input.DeviceProfileID
		profileID = &val
	}
	if input.CabinetID != "" {
		val := input.CabinetID
		cabinetID = &val
	}
	if input.TeamID != "" {
		val := input.TeamID
		teamID = &val
	}

	device := &domain.Device{
		Name:            input.Name,
		Description:     input.Description,
		SerialNumber:    input.SerialNumber,
		EPC:             input.EPC,
		DeviceProfileID: profileID,
		CabinetID:       cabinetID,
		TeamID:          teamID,
		TenantID:        input.TenantID,
	}

	if err := uc.svc.Create(ctx, device); err != nil {
		return nil, err
	}

	return &DeviceOutput{
		ID:           device.ID,
		Name:         device.Name,
		SerialNumber: device.SerialNumber,
		Description:  device.Description,
		EPC:          device.EPC,
	}, nil
}

func (uc *DeviceUseCase) GetByID(ctx context.Context, id string) (*DeviceOutput, error) {
	device, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &DeviceOutput{
		ID:           device.ID,
		Name:         device.Name,
		SerialNumber: device.SerialNumber,
		Description:  device.Description,
		EPC:          device.EPC,
	}, nil
}

func (uc *DeviceUseCase) List(ctx context.Context, limit, offset int) ([]DeviceOutput, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}
	devices, err := uc.svc.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	output := make([]DeviceOutput, len(devices))
	for i, d := range devices {
		output[i] = DeviceOutput{
			ID:           d.ID,
			Name:         d.Name,
			SerialNumber: d.SerialNumber,
			Description:  d.Description,
			EPC:          d.EPC,
		}
	}
	return output, nil
}

func (uc *DeviceUseCase) Update(ctx context.Context, id string, input CreateDeviceInput) (*DeviceOutput, error) {
	device, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	device.Name = input.Name
	device.Description = input.Description
	device.EPC = input.EPC
	if input.SerialNumber != "" {
		device.SerialNumber = input.SerialNumber
	}

	if err := uc.svc.Update(ctx, device); err != nil {
		return nil, err
	}

	return &DeviceOutput{
		ID:           device.ID,
		Name:         device.Name,
		SerialNumber: device.SerialNumber,
		Description:  device.Description,
		EPC:          device.EPC,
	}, nil
}

func (uc *DeviceUseCase) Delete(ctx context.Context, id string) error {
	return uc.svc.Delete(ctx, id)
}
