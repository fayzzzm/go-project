package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlerFunc matches the signature of a controller method that returns data or error.
type HandlerFunc func(c *gin.Context) (any, error)

// Handle wraps a controller function to standardize response and error handling.
// It executes the function, handles any returned error via Gin's error mechanism,
// or sends a JSON response with the successStatus.
func Handle(fn HandlerFunc, successStatus int) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := fn(c)
		if err != nil {
			c.Error(err)
			return
		}

		if successStatus == http.StatusNoContent {
			c.Status(http.StatusNoContent)
			return
		}

		if res == nil {
			// If we expect content but got nil, maybe send empty object or just status?
			// Sending generic map for empty JSON if 200 OK.
			// Or allow nil -> null in JSON.
			// For consistency, let standard JSON marshaler handle nil (it outputs "null").
		}

		c.JSON(successStatus, res)
	}
}
