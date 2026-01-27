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
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description"`
	SerialNumber    string `json:"serial_number" binding:"required"`
	EPC             string `json:"epc"`
	DeviceProfileID string `json:"device_profile_id"`
	CabinetID       string `json:"cabinet_id"`
	TeamID          string `json:"team_id"`
	TenantID        string `json:"tenant_id"`
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

	// Handle Description and EPC pointers
	var desc, epc *string
	if input.Description != "" {
		val := input.Description
		desc = &val
	}
	if input.EPC != "" {
		val := input.EPC
		epc = &val
	}

	var tenantID *string
	if input.TenantID != "" {
		val := input.TenantID
		tenantID = &val
	}

	device := &domain.Device{
		Name:            input.Name,
		Description:     desc,
		SerialNumber:    input.SerialNumber,
		EPC:             epc,
		DeviceProfileID: profileID,
		CabinetID:       cabinetID,
		TeamID:          teamID,
		TenantID:        tenantID,
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
	limit := p.Limit
	offset := p.Offset
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
		output[i] = *toDeviceOutput(&d)
	}
	return output, nil
}

func (uc *DeviceUseCase) Update(ctx context.Context, id string, input CreateDeviceInput) (*DeviceOutput, error) {
	device, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Security Check: Ensure Device belongs to the requested Tenant
	if input.TenantID != "" {
		if device.TenantID != nil && *device.TenantID != input.TenantID {
			return nil, domain.ErrUnauthorized
		}
	}

	device.Name = input.Name
	if input.Description != "" {
		val := input.Description
		device.Description = &val
	}
	if input.EPC != "" {
		val := input.EPC
		device.EPC = &val
	}
	if input.SerialNumber != "" {
		device.SerialNumber = input.SerialNumber
	}

	if input.DeviceProfileID != "" {
		val := input.DeviceProfileID
		device.DeviceProfileID = &val
	}
	if input.CabinetID != "" {
		val := input.CabinetID
		device.CabinetID = &val
	}
	if input.TeamID != "" {
		val := input.TeamID
		device.TeamID = &val
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

func (uc *DeviceUseCase) Delete(ctx context.Context, id, tenantID string) error {
	// Security Check
	if tenantID != "" {
		device, err := uc.svc.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if device.TenantID != nil && *device.TenantID != tenantID {
			return domain.ErrUnauthorized
		}
	}
	return uc.svc.Delete(ctx, id)
}
