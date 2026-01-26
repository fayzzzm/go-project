package domain

import (
	"github.com/fayzzzm/go-project/pkg/validator"
)

func (d *Device) Validate() (bool, string) {
	base := validator.BaseValidator{}

	if ok, msg := base.StringValidator(d.Name, 3, 100, true); !ok {
		return false, "Name: " + msg
	}
	if ok, msg := base.StringValidator(d.SerialNumber, 5, 50, true); !ok {
		return false, "SerialNumber: " + msg
	}

	if d.TenantID == "" {
		return false, "TenantID is required"
	}

	return true, ""
}

func (dp *DeviceProfile) Validate() (bool, string) {
	base := validator.BaseValidator{}
	if ok, msg := base.StringValidator(dp.Name, 1, 100, true); !ok {
		return false, "Name: " + msg
	}
	return true, ""
}

func (u *User) Validate() (bool, string) {
	emailVal := validator.NewEmailValidator()
	if ok, msg := emailVal.Validate(u.Email); !ok {
		return false, "Email: " + msg
	}
	return true, ""
}
