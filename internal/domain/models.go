package domain

import (
	"time"
)

type Pagination struct {
	Limit  int
	Offset int
}

func (p *Pagination) Normalize() {
	if p.Limit <= 0 {
		p.Limit = 100
	}
	if p.Limit > 1000 {
		p.Limit = 1000
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

// Tenant represents a customer or organization.
type Tenant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// DeviceProfile represents metadata about a type of device.
type DeviceProfile struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	TenantID    *string   `json:"tenant_id" db:"tenant_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Member struct {
	TeamID    string    `json:"team_id" db:"team_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Device represents a specific instance of a device profile.
type Device struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     *string   `json:"description" db:"description"`
	SerialNumber    string    `json:"serial_number" db:"serial_number"`
	EPC             *string   `json:"epc" db:"epc"`
	DeviceProfileID *string   `json:"device_profile_id" db:"device_profile_id"`
	Status          string    `json:"status" db:"status"`
	CabinetID       *string   `json:"cabinet_id" db:"cabinet_id"`
	TeamID          *string   `json:"team_id" db:"team_id"`
	TenantID        *string   `json:"tenant_id" db:"tenant_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Cabinet represents a storage cabinet.
type Cabinet struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Location    *string   `json:"location" db:"location"` // Address/Location
	MachineID   *string   `json:"machine_id" db:"machine_id"`
	Status      *string   `json:"status" db:"status"`
	TeamID      *string   `json:"team_id" db:"team_id"`
	TenantID    *string   `json:"tenant_id" db:"tenant_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy   string    `json:"updated_by" db:"updated_by"`
	Team        *Team     `json:"team,omitempty" db:"-"` // Association, ignore in DB scan if flat
}

// Team represents a group of users or context.
type Team struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Status    *string   `json:"status" db:"status"`
	TenantID  *string   `json:"tenant_id" db:"tenant_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy string    `json:"updated_by" db:"updated_by"`
}

// User represents a system user.
type User struct {
	ID                 string                 `json:"id" db:"id"`
	Email              string                 `json:"email" db:"email"`
	Phone              *string                `json:"phone" db:"phone"`
	Role               string                 `json:"role" db:"role"`
	TenantID           string                 `json:"tenant_id" db:"tenant_id"`
	Aud                string                 `json:"aud" db:"-"`
	IsAnonymous        bool                   `json:"is_anonymous" db:"-"`
	ConfirmationSentAt time.Time              `json:"confirmation_sent_at" db:"-"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`
	AppMetadata        map[string]interface{} `json:"app_metadata" db:"app_metadata"`   // JSONB
	UserMetadata       map[string]interface{} `json:"user_metadata" db:"user_metadata"` // JSONB
	// Name field requested by user, mapping to user_metadata potentially or top level if extending schema
	Name     *string `json:"name,omitempty" db:"name"`
	Address  *string `json:"address,omitempty" db:"address"`
	Password string  `json:"-" db:"password"`
}
