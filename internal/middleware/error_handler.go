package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

// APIErrorResponse represents a standard error response structure.
type APIErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationErrorDetail represents a single field validation error.
type ValidationErrorDetail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// ErrorHandlerMiddleware intercepts errors attached to the context and sends structured JSON responses.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Execute the handlers first

		// After handlers return, check if there are errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err
			var parsedError *APIErrorResponse
			statusCode := http.StatusInternalServerError

			// 1. Check for Validation Errors
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				statusCode = http.StatusBadRequest
				details := make([]ValidationErrorDetail, len(ve))
				for i, fe := range ve {
					details[i] = ValidationErrorDetail{
						Field:  utils.ToSnakeCase(fe.Field()), // Ensure snake_case for JSON fields
						Reason: msgForTag(fe),
					}
				}
				parsedError = &APIErrorResponse{
					Code:    domain.CodeValidationFailed,
					Message: "Input validation failed",
					Details: details,
				}
			} else if unmarshalTypeError := (*json.UnmarshalTypeError)(nil); errors.As(err, &unmarshalTypeError) {
				statusCode = http.StatusBadRequest
				parsedError = &APIErrorResponse{
					Code:    domain.CodeInvalidJSON,
					Message: fmt.Sprintf("Field '%s' expected type '%v' but got '%v'", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value),
				}
			} else {
				// 2. Check Domain Concept Errors
				switch {
				case errors.Is(err, domain.ErrNotFound):
					statusCode = http.StatusNotFound
					parsedError = &APIErrorResponse{Code: domain.CodeNotFound, Message: err.Error()}
				case errors.Is(err, domain.ErrConflict):
					statusCode = http.StatusConflict
					parsedError = &APIErrorResponse{Code: domain.CodeConflict, Message: err.Error()}
				case errors.Is(err, domain.ErrInvalidInput):
					statusCode = http.StatusBadRequest
					parsedError = &APIErrorResponse{Code: domain.CodeInvalidInput, Message: err.Error()}
				case errors.Is(err, domain.ErrUnauthorized):
					statusCode = http.StatusUnauthorized
					parsedError = &APIErrorResponse{Code: domain.CodeUnauthorized, Message: err.Error()}
				case errors.Is(err, domain.ErrForbidden):
					statusCode = http.StatusForbidden
					parsedError = &APIErrorResponse{Code: domain.CodeForbidden, Message: err.Error()}
				case errors.Is(err, domain.ErrInvalidCredentials):
					statusCode = http.StatusUnauthorized
					parsedError = &APIErrorResponse{Code: domain.CodeInvalidCredentials, Message: err.Error()}
				default:
					// 3. Handle Postgres specific errors
					var pgErr *pgconn.PgError
					if errors.As(err, &pgErr) {
						switch pgErr.Code {
						case domain.PGUniqueViolation:
							statusCode = http.StatusConflict
							parsedError = &APIErrorResponse{
								Code:    domain.CodeResourceExists,
								Message: "A resource with these unique identifiers already exists.",
							}
						case domain.PGForeignKeyViolation:
							statusCode = http.StatusBadRequest
							parsedError = &APIErrorResponse{
								Code:    domain.CodeConstraintViolation,
								Message: "Operation violates foreign key constraints (referenced resource may not exist).",
							}
						case domain.PGCustomException: // Custom Raise Exception
							statusCode = http.StatusForbidden
							parsedError = &APIErrorResponse{
								Code:    domain.CodeBusinessRuleViolation,
								Message: pgErr.Message,
							}
						default:
							// Log unexpected DB errors
							fmt.Printf("[DB Error] %s\n", pgErr.Error())
							parsedError = &APIErrorResponse{
								Code:    domain.CodeInternalError,
								Message: "An internal database error occurred.",
							}
						}
					} else {
						// 4. Default Internal Server Error
						fmt.Printf("[Internal Error] %s\n", err.Error())
						parsedError = &APIErrorResponse{
							Code:    domain.CodeInternalError,
							Message: fmt.Sprintf("Internal Error: %s", err.Error()),
						}
					}
				}
			}

			// If status is already written, we can't do anything
			if !c.Writer.Written() {
				c.JSON(statusCode, gin.H{"error": parsedError})
			}
		}
	}
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case domain.ValidationTagRequired:
		return "This field is required"
	case domain.ValidationTagEmail:
		return "Invalid email format"
	case domain.ValidationTagMin:
		return fmt.Sprintf("Must be at least %s characters long", fe.Param())
	case domain.ValidationTagUUID:
		return "Must be a valid UUID"
	default:
		return fmt.Sprintf("Failed validation on tag '%s'", fe.Tag())
	}
}
