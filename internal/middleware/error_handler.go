package middleware

import (
	"errors"
	"net/http"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/gin-gonic/gin"
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
			default:
				statusCode = http.StatusInternalServerError
				message = "internal server error"
				// In production, log the real error here but keep the response generic
			}

			// If status is already written, we can't do anything (rare)
			if !c.Writer.Written() {
				c.JSON(statusCode, gin.H{"error": message})
			}
		}
	}
}
