package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Email    string `json:"email_address" validate:"required,email"`
	Age      int    `json:"user_age" validate:"min=18"`
	ID       string `json:"id" validate:"uuid"`
	Username string `json:"username" validate:"required"`
}

func TestValidate(t *testing.T) {
	t.Run("Valid Struct", func(t *testing.T) {
		s := testStruct{
			Email:    "test@example.com",
			Age:      20,
			ID:       "550e8400-e29b-41d4-a716-446655440000",
			Username: "johndoe",
		}
		err := Validate(s)
		assert.NoError(t, err)
	})

	t.Run("Required Fields Missing", func(t *testing.T) {
		s := testStruct{}
		err := Validate(s)
		assert.Error(t, err)

		errors := GetErrors(err)
		assert.NotEmpty(t, errors)

		msgs := make(map[string]string)
		for _, e := range errors {
			msgs[e.Field] = e.Reason
		}

		assert.Equal(t, "This field is required", msgs["email_address"])
		assert.Equal(t, "This field is required", msgs["username"])
		assert.Equal(t, "Must be at least 18", msgs["user_age"])
		assert.Equal(t, "Must be a valid UUID", msgs["id"])
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		s := testStruct{Email: "invalid-email", Age: 20, Username: "test", ID: "550e8400-e29b-41d4-a716-446655440000"}
		err := Validate(s)
		errors := GetErrors(err)

		var emailErr *ValidationErrorDetail
		for _, e := range errors {
			if e.Field == "email_address" {
				emailErr = &e
				break
			}
		}

		assert.NotNil(t, emailErr)
		assert.Equal(t, "Invalid email format", emailErr.Reason)
	})

	t.Run("Min Value Constraint", func(t *testing.T) {
		s := testStruct{Email: "test@example.com", Age: 10, Username: "test", ID: "550e8400-e29b-41d4-a716-446655440000"}
		err := Validate(s)
		errors := GetErrors(err)

		var ageErr *ValidationErrorDetail
		for _, e := range errors {
			if e.Field == "user_age" {
				ageErr = &e
				break
			}
		}

		assert.NotNil(t, ageErr)
		assert.Equal(t, "Must be at least 18", ageErr.Reason)
	})

	t.Run("Invalid UUID Format", func(t *testing.T) {
		s := testStruct{ID: "not-a-uuid"}
		err := Validate(s)
		errors := GetErrors(err)

		var idErr *ValidationErrorDetail
		for _, e := range errors {
			if e.Field == "id" {
				idErr = &e
				break
			}
		}

		assert.NotNil(t, idErr)
		assert.Equal(t, "Must be a valid UUID", idErr.Reason)
	})
}
