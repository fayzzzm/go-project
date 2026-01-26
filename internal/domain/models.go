package domain

import (
	"time"
)

// DeviceProfile represents metadata about a type of device.
type DeviceProfile struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	TenantID    *string   `json:"tenant_id" db:"tenant_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Device represents a specific instance of a device profile.
type Device struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	SerialNumber     string    `json:"serial_number" db:"serial_number"`
	EPC              string    `json:"epc" db:"epc"`
	DeviceProfileID  *string   `json:"device_profile_id" db:"device_profile_id"`
	DeviceStatusID   string    `json:"device_status_id" db:"device_status_id"`
	CabinetID        *string   `json:"cabinet_id" db:"cabinet_id"`
	TeamID           *string   `json:"team_id" db:"team_id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	Owner            string    `json:"owner" db:"owner"`
	OwnerModifiedAt  time.Time `json:"owner_modified_at" db:"owner_modified_at"`
	OwnerModifiedBy  string    `json:"owner_modified_by" db:"owner_modified_by"`
	LastCheckedOutBy string    `json:"last_checked_out_by" db:"last_checked_out_by"`
	LastModifiedBy   string    `json:"last_modified_by" db:"last_modified_by"`
	LastModifiedOn   time.Time `json:"last_modified_on" db:"last_modified_on"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	CreatedBy        string    `json:"created_by" db:"created_by"`
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
	Aud                string                 `json:"aud" db:"-"`
	IsAnonymous        bool                   `json:"is_anonymous" db:"-"`
	ConfirmationSentAt time.Time              `json:"confirmation_sent_at" db:"-"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`
	AppMetadata        map[string]interface{} `json:"app_metadata" db:"app_metadata"`   // JSONB
	UserMetadata       map[string]interface{} `json:"user_metadata" db:"user_metadata"` // JSONB
	// Name field requested by user, mapping to user_metadata potentially or top level if extending schema
	Name    *string `json:"name,omitempty" db:"name"`
	Address *string `json:"address,omitempty" db:"address"`
}
