package middleware

import (
	"strconv"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/gin-gonic/gin"
)

// GetPagination extracts limit and offset from the query parameters.
func GetPagination(c *gin.Context) domain.Pagination {
	var p domain.Pagination
	if l := c.Query("limit"); l != "" {
		p.Limit, _ = strconv.Atoi(l)
	}
	if o := c.Query("offset"); o != "" {
		p.Offset, _ = strconv.Atoi(o)
	}
	return p
}
