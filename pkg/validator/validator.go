package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// ValidationResult represents the result of a validation.
type ValidationResult struct {
	IsValid bool
	Message string
}

// ValidatorFunc defines the signature for validation functions.
type ValidatorFunc func(interface{}) (bool, string)

// BaseValidator provides basic validation methods.
type BaseValidator struct{}

// StringValidator validates string constraints.
func (v BaseValidator) StringValidator(val string, minLen, maxLen int, required bool) (bool, string) {
	if required && val == "" {
		return false, "value is required"
	}
	if !required && val == "" {
		return true, ""
	}
	length := utf8.RuneCountInString(val)
	if length < minLen || length > maxLen {
		return false, "value length must be between " + string(rune(minLen)) + " and " + string(rune(maxLen))
	}
	return true, ""
}

// EmailValidator validates email format.
type EmailValidator struct {
	BaseValidator
	MinLen int
	MaxLen int
}

func NewEmailValidator() *EmailValidator {
	return &EmailValidator{
		MinLen: 5,
		MaxLen: 254,
	}
}

func (v *EmailValidator) Validate(email string) (bool, string) {
	// First check basic string properties using BaseValidator
	if ok, msg := v.StringValidator(email, v.MinLen, v.MaxLen, true); !ok {
		return false, msg
	}

	// Then check specific email format
	// Simple regex for demonstration
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(strings.ToLower(email)) {
		return false, "invalid email format"
	}
	return true, ""
}

// ValidationEngine handles running validations on structs.
type ValidationEngine struct{}

func NewValidationEngine() *ValidationEngine {
	return &ValidationEngine{}
}
