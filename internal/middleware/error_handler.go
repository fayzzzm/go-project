package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/validate"
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

// ErrorHandlerMiddleware intercepts errors attached to the context and sends structured JSON responses.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	v := validate.NewValidator()

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var parsedError *APIErrorResponse
			statusCode := http.StatusInternalServerError

			// 1. Handle Domain Errors (Highly Professional)
			var de *domain.DomainError
			if errors.As(err, &de) {
				statusCode = de.HTTPStatus
				parsedError = &APIErrorResponse{Code: de.Code, Message: de.Message}
			} else if ve, ok := err.(validator.ValidationErrors); ok || errors.As(err, &ve) {
				// 2. Handle Validation Errors
				statusCode = http.StatusBadRequest
				parsedError = &APIErrorResponse{
					Code:    domain.CodeValidationFailed,
					Message: "Input validation failed",
					Details: v.GetErrors(err),
				}
			} else if unmarshalTypeError := (*json.UnmarshalTypeError)(nil); errors.As(err, &unmarshalTypeError) {
				// 3. Handle JSON Syntax/Type Errors
				statusCode = http.StatusBadRequest
				parsedError = &APIErrorResponse{
					Code:    domain.CodeInvalidJSON,
					Message: fmt.Sprintf("Field '%s' expected type '%v' but got '%v'", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value),
				}
			} else if pgErr := (*pgconn.PgError)(nil); errors.As(err, &pgErr) {
				// 4. Handle DB specific errors (Optimized)
				statusCode, parsedError = handlePostgresError(pgErr)
			} else {
				// 5. Default Fallback
				fmt.Printf("[Internal Error] %s\n", err.Error())
				parsedError = &APIErrorResponse{
					Code:    domain.CodeInternalError,
					Message: "An unexpected error occurred.",
				}
			}

			if !c.Writer.Written() {
				c.JSON(statusCode, gin.H{"error": parsedError})
			}
		}
	}
}

func handlePostgresError(pgErr *pgconn.PgError) (int, *APIErrorResponse) {
	switch pgErr.Code {
	case domain.PGUniqueViolation:
		return http.StatusConflict, &APIErrorResponse{Code: domain.CodeResourceExists, Message: "Resource already exists."}
	case domain.PGForeignKeyViolation:
		return http.StatusBadRequest, &APIErrorResponse{Code: domain.CodeConstraintViolation, Message: "Constraint violation."}
	case domain.PGCustomException:
		return http.StatusForbidden, &APIErrorResponse{Code: domain.CodeBusinessRuleViolation, Message: pgErr.Message}
	default:
		fmt.Printf("[DB Error] %s\n", pgErr.Error())
		return http.StatusInternalServerError, &APIErrorResponse{Code: domain.CodeInternalError, Message: "Database error."}
	}
}
