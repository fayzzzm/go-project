package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handle wraps a controller function to standardize response and error handling.
func Handle[T any](fn func(c *gin.Context) (T, error), successStatus int) gin.HandlerFunc {
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

		c.JSON(successStatus, res)
	}
}
