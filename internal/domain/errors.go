package domain

import "errors"

const (
	// Error Codes
	CodeValidationFailed      = "VALIDATION_FAILED"
	CodeInvalidJSON           = "INVALID_JSON_TYPE"
	CodeNotFound              = "NOT_FOUND"
	CodeConflict              = "CONFLICT"
	CodeInvalidInput          = "INVALID_INPUT"
	CodeUnauthorized          = "UNAUTHORIZED"
	CodeForbidden             = "FORBIDDEN"
	CodeInvalidCredentials    = "INVALID_CREDENTIALS"
	CodeResourceExists        = "RESOURCE_EXISTS"
	CodeConstraintViolation   = "CONSTRAINT_VIOLATION"
	CodeBusinessRuleViolation = "BUSINESS_RULE_VIOLATION"
	CodeInternalError         = "INTERNAL_ERROR"

	// Validation Tags (from go-playground/validator)
	ValidationTagRequired = "required"
	ValidationTagEmail    = "email"
	ValidationTagMin      = "min"
	ValidationTagUUID     = "uuid"

	// Postgres Error Codes
	PGUniqueViolation     = "23505"
	PGForeignKeyViolation = "23503"
	PGCustomException     = "P0001" // Used by RAISE EXCEPTION in PL/pgSQL
)

var (
	// Standard Errors
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource already exists") // duplicate key
	ErrInternal           = errors.New("internal system error")
	ErrInvalidInput       = errors.New("invalid input parameter")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("access forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
