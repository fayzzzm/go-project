package domain

import (
	"github.com/fayzzzm/go-project/pkg/validator"
)

// Validate implements the Validatable interface for Device.
func (d *Device) Validate() (bool, string) {
	base := validator.BaseValidator{}

	// Validate Name
	if ok, msg := base.StringValidator(d.Name, 3, 100, true); !ok {
		return false, "Name: " + msg
	}

	// Validate SerialNumber
	if ok, msg := base.StringValidator(d.SerialNumber, 5, 50, true); !ok {
		return false, "SerialNumber: " + msg
	}

	// Validate TenantID (example of required field)
	if d.TenantID == "" {
		return false, "TenantID is required"
	}

	return true, ""
}

// Validate implements the Validatable interface for DeviceProfile.
func (dp *DeviceProfile) Validate() (bool, string) {
	base := validator.BaseValidator{}
	if ok, msg := base.StringValidator(dp.Name, 1, 100, true); !ok {
		return false, "Name: " + msg
	}
	return true, ""
}

// Validate implements the Validatable interface for User.
func (u *User) Validate() (bool, string) {
	emailVal := validator.NewEmailValidator()
	if ok, msg := emailVal.Validate(u.Email); !ok {
		return false, "Email: " + msg
	}
	return true, ""
}
