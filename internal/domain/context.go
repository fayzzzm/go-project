package domain

import "context"

type contextKey string

const (
	tenantKey contextKey = "tenant_id"
	userKey   contextKey = "user_id"
)

// ContextWithTenant returns a new context with the tenant ID
func ContextWithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

// TenantFromContext returns the tenant ID from the context
func TenantFromContext(ctx context.Context) string {
	val, _ := ctx.Value(tenantKey).(string)
	return val
}

// ContextWithUser returns a new context with the user ID
func ContextWithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

// UserFromContext returns the user ID from the context
func UserFromContext(ctx context.Context) string {
	val, _ := ctx.Value(userKey).(string)
	return val
}
