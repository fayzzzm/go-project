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
		user, _ := c.Get(ContextUserID) // Assume set by Auth
		tenant, _ := c.Get(ContextTenantID)

		// TODO: Connect a NoSQL database (e.g., MongoDB, ElasticSearch) here to persist audit logs.
		// For SOC 2 compliance, logs should be retained securely and be searchable.
		// Example: mongoClient.Database("audit").Collection("logs").InsertOne(ctx, auditEntry)
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
