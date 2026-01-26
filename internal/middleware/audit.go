package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditMiddleware logs requests for SOC 2 compliance.
// It records who, when, and what was accessed.
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		user, _ := c.Get("user_id") // Assume set by Auth
		tenant, _ := c.Get("tenant_id")

		log.Printf("[AUDIT] %s | %d | %s | %s | User: %v | Tenant: %v | Latency: %v",
			start.Format(time.RFC3339),
			status,
			method,
			path,
			user,
			tenant,
			latency,
		)
	}
}
