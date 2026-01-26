package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPagination extracts limit and offset from the query parameters.
// It returns (0, 0) if parsing fails or parameters are missing,
// delegating default logic to the business layer.
func GetPagination(c *gin.Context) (int, int) {
	limit := 0
	if l := c.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	offset := 0
	if o := c.Query("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}
	return limit, offset
}
