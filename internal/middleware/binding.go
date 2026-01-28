package middleware

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// BodyContextKey is the key used to store the parsed body in the context.
const BodyContextKey = "request_body"

// BindJSON is a generic middleware that binds the request body to a struct of type T.
func BindJSON[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input T
		if err := c.ShouldBindJSON(&input); err != nil {
			// Pass raw error to global error handler to format validation messages
			_ = c.Error(err).SetType(gin.ErrorTypeBind)
			c.Abort()
			return
		}

		// Inject TenantID: Check Context (Auth) first, then Header
		tenantID := c.GetString(ContextTenantID)
		if tenantID == "" {
			tenantID = c.GetHeader("X-Tenant-ID")
		}

		if tenantID != "" {
			val := reflect.ValueOf(&input).Elem()
			if val.Kind() == reflect.Struct {
				field := val.FieldByName("TenantID")
				if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
					if field.String() == "" {
						field.SetString(tenantID)
					}
				}
			}
		}

		c.Set(BodyContextKey, input)
		c.Next()
	}
}

// GetBody retrieves the bound body from the context and casts it to type T.
func GetBody[T any](c *gin.Context) T {
	val, exists := c.Get(BodyContextKey)
	if !exists {
		panic("GetBody called on handler without BindJSON middleware")
	}
	return val.(T)
}
