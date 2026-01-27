package postgres

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
	LimitVal  *int    `db:"limit_val"`
	OffsetVal *int    `db:"offset_val"`
}
