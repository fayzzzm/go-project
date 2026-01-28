package postgres

import (
	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
)

type UserRequest struct {
	ID           *string                `db:"id"`
	Email        *string                `db:"email"`
	Name         *string                `db:"name"`
	Address      *string                `db:"address"`
	Phone        *string                `db:"phone"`
	Role         *string                `db:"role"`
	AppMetadata  map[string]interface{} `db:"app_metadata"`
	UserMetadata map[string]interface{} `db:"user_metadata"`
	Password     *string                `db:"password"`
	LimitVal     *int                   `db:"limit_val"`
	OffsetVal    *int                   `db:"offset_val"`
	TeamID       *string                `db:"team_id"`
	TenantID     *string                `db:"tenant_id"`
}

type DeviceRequest struct {
	ID              *string `db:"id"`
	Name            *string `db:"name"`
	Description     *string `db:"description"`
	SerialNumber    *string `db:"serial_number"`
	EPC             *string `db:"epc"`
	DeviceProfileID *string `db:"device_profile_id"`
	CabinetID       *string `db:"cabinet_id"`
	TeamID          *string `db:"team_id"`
	TenantID        *string `db:"tenant_id"`
	LimitVal        *int    `db:"limit_val"`
	OffsetVal       *int    `db:"offset_val"`
}

type CabinetRequest struct {
	ID          *string `db:"id"`
	Name        *string `db:"name"`
	Description *string `db:"description"`
	Location    *string `db:"location"`
	MachineID   *string `db:"machine_id"`
	Status      *string `db:"status"`
	TeamID      *string `db:"team_id"`
	TenantID    *string `db:"tenant_id"`
	LimitVal    *int    `db:"limit_val"`
	OffsetVal   *int    `db:"offset_val"`
	UserID      *string `db:"user_id"`
}

type TeamRequest struct {
	ID        *string `db:"id"`
	Name      *string `db:"name"`
	Status    *string `db:"status"`
	TenantID  *string `db:"tenant_id"`
	LimitVal  *int    `db:"limit_val"`
	OffsetVal *int    `db:"offset_val"`
	UserID    *string `db:"user_id"`
}

type DeviceProfileRequest struct {
	ID          *string `db:"id"`
	Name        *string `db:"name"`
	Description *string `db:"description"`
	TenantID    *string `db:"tenant_id"`
	LimitVal    *int    `db:"limit_val"`
	OffsetVal   *int    `db:"offset_val"`
}

type MemberRequest struct {
	TeamID *string `db:"team_id"`
	UserID *string `db:"user_id"`
	Role   *string `db:"role"`
}

type TenantRequest struct {
	ID        *string `db:"id"`
	Name      *string `db:"name"`
	LimitVal  *int    `db:"limit_val"`
	OffsetVal *int    `db:"offset_val"`
}

func NewUserRequest(u *domain.User) *UserRequest {
	return &UserRequest{
		ID:           utils.StringPtrOrNil(u.ID),
		Email:        utils.StringPtrOrNil(u.Email),
		Name:         u.Name,
		Address:      u.Address,
		Phone:        u.Phone,
		Role:         utils.StringPtrOrNil(u.Role),
		AppMetadata:  u.AppMetadata,
		UserMetadata: u.UserMetadata,
		Password:     utils.StringPtrOrNil(u.Password),
		TenantID:     utils.StringPtrOrNil(u.TenantID),
	}
}

func NewTeamRequest(t *domain.Team) *TeamRequest {
	return &TeamRequest{
		ID:       utils.StringPtrOrNil(t.ID),
		Name:     &t.Name,
		Status:   t.Status,
		TenantID: t.TenantID,
	}
}

func NewDeviceRequest(d *domain.Device) *DeviceRequest {
	return &DeviceRequest{
		ID:              utils.StringPtrOrNil(d.ID),
		Name:            &d.Name,
		Description:     d.Description,
		SerialNumber:    &d.SerialNumber,
		EPC:             d.EPC,
		DeviceProfileID: d.DeviceProfileID,
		CabinetID:       d.CabinetID,
		TeamID:          d.TeamID,
		TenantID:        d.TenantID,
	}
}

func NewCabinetRequest(c *domain.Cabinet, userID *string) *CabinetRequest {
	return &CabinetRequest{
		ID:          utils.StringPtrOrNil(c.ID),
		Name:        &c.Name,
		Description: c.Description,
		Location:    c.Location,
		MachineID:   c.MachineID,
		Status:      c.Status,
		TeamID:      c.TeamID,
		TenantID:    c.TenantID,
		UserID:      userID,
	}
}

func NewDeviceProfileRequest(dp *domain.DeviceProfile) *DeviceProfileRequest {
	return &DeviceProfileRequest{
		ID:          utils.StringPtrOrNil(dp.ID),
		Name:        &dp.Name,
		Description: dp.Description,
		TenantID:    dp.TenantID,
	}
}
