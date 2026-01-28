package domain

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

// DomainError is a structured error with a machine-readable code
type DomainError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *DomainError) Error() string { return e.Message }

var (
	ErrNotFound           = &DomainError{Code: CodeNotFound, Message: "resource not found", HTTPStatus: 404}
	ErrConflict           = &DomainError{Code: CodeConflict, Message: "resource already exists", HTTPStatus: 409}
	ErrInvalidInput       = &DomainError{Code: CodeInvalidInput, Message: "invalid input parameter", HTTPStatus: 400}
	ErrUnauthorized       = &DomainError{Code: CodeUnauthorized, Message: "unauthorized access", HTTPStatus: 401}
	ErrForbidden          = &DomainError{Code: CodeForbidden, Message: "access forbidden", HTTPStatus: 403}
	ErrInvalidCredentials = &DomainError{Code: CodeInvalidCredentials, Message: "invalid credentials", HTTPStatus: 401}
	ErrInternal           = &DomainError{Code: CodeInternalError, Message: "internal system error", HTTPStatus: 500}
)
