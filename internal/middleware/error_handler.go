package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorHandlerMiddleware intercepts errors attached to the context and sends generic JSON responses.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Execute the handlers first

		// After handlers return, check if there are errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err

			var statusCode int
			var message string

			// Inspect for Validation Errors first (specific types)
			var ve validator.ValidationErrors
			var unmarshalTypeError *json.UnmarshalTypeError

			if errors.As(err, &ve) {
				statusCode = http.StatusBadRequest
				out := make([]string, len(ve))
				for i, fe := range ve {
					out[i] = fmt.Sprintf("Field '%s' failed validation (tag: %s)", fe.Field(), fe.Tag())
				}
				message = strings.Join(out, "; ")
			} else if errors.As(err, &unmarshalTypeError) {
				statusCode = http.StatusBadRequest
				message = fmt.Sprintf("Field '%s' expected type '%v' but got '%v'", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)
			} else {
				// Map Domain Errors to HTTP Status Codes
				switch {
				case errors.Is(err, domain.ErrNotFound):

					statusCode = http.StatusNotFound
					message = err.Error()
				case errors.Is(err, domain.ErrConflict):
					statusCode = http.StatusConflict
					message = err.Error()
				case errors.Is(err, domain.ErrInvalidInput):
					statusCode = http.StatusBadRequest
					message = err.Error()
				case errors.Is(err, domain.ErrUnauthorized):
					statusCode = http.StatusUnauthorized
					message = err.Error()
				case errors.Is(err, domain.ErrForbidden):
					statusCode = http.StatusForbidden
					message = err.Error()
				case errors.Is(err, domain.ErrInvalidCredentials):
					statusCode = http.StatusUnauthorized
					message = err.Error()
				default:
					// Handle Postgres custom errors
					var pgErr *pgconn.PgError
					if errors.As(err, &pgErr) {
						if pgErr.Code == "P0001" {
							statusCode = http.StatusForbidden
							message = pgErr.Message
						} else {
							// Log other DB errors
							println("DB Error:", pgErr.Error())
							statusCode = http.StatusInternalServerError
							message = "internal server error"
						}
					} else {
						// Log the actual error for debugging
						// TODO: Use a proper logger
						println("Internal Server Error:", err.Error())
						statusCode = http.StatusInternalServerError
						message = "internal server error"
					}
				} // End switch
			} // End if/else validation checks

			// If status is already written, we can't do anything (rare)
			if !c.Writer.Written() {
				c.JSON(statusCode, gin.H{"error": message})
			}
		}
	}
}
