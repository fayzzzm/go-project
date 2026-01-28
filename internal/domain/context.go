package domain

import "context"

type contextKey string

const tenantKey contextKey = "tenant_id"

// ContextWithTenant returns a new context with the tenant ID
func ContextWithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

// TenantFromContext returns the tenant ID from the context
func TenantFromContext(ctx context.Context) string {
	val, _ := ctx.Value(tenantKey).(string)
	return val
}
