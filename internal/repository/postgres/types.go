package postgres

// DeviceRequest matches a PostgreSQL composite type for device operations
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

// UserRequest matches a PostgreSQL composite type for user operations
type UserRequest struct {
	ID           *string                `db:"id"`
	Email        *string                `db:"email"`
	Name         *string                `db:"name"`
	Address      *string                `db:"address"`
	Phone        *string                `db:"phone"`
	Role         *string                `db:"role"`
	AppMetadata  map[string]interface{} `db:"app_metadata"`
	UserMetadata map[string]interface{} `db:"user_metadata"`
	LimitVal     *int                   `db:"limit_val"`
	OffsetVal    *int                   `db:"offset_val"`
}

// CabinetRequest matches a PostgreSQL composite type for cabinet operations
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
}

// TeamRequest matches a PostgreSQL composite type for team operations
type TeamRequest struct {
	ID        *string `db:"id"`
	Name      *string `db:"name"`
	Status    *string `db:"status"`
	TenantID  *string `db:"tenant_id"`
	LimitVal  *int    `db:"limit_val"`
	OffsetVal *int    `db:"offset_val"`
}

// DeviceProfileRequest matches a PostgreSQL composite type for device profile operations
type DeviceProfileRequest struct {
	ID          *string `db:"id"`
	Name        *string `db:"name"`
	Description *string `db:"description"`
	TenantID    *string `db:"tenant_id"`
	LimitVal    *int    `db:"limit_val"`
	OffsetVal   *int    `db:"offset_val"`
}
