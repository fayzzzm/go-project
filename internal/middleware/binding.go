package middleware

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// BodyContextKey is the key used to store the parsed body in the context.
const BodyContextKey = "request_body"

// BindJSON is a generic middleware that binds the request body to a struct of type T.
func BindJSON[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input T
		// 1. Initial binding from JSON body
		if err := c.ShouldBindJSON(&input); err != nil {
			// We ignore certain errors here if they might be fixed by injection,
			// but it's cleaner to just bind and then validate.
			// Actually Gin's ShouldBindJSON runs validator immediately.
		}

		// 2. Inject TenantID from context (set by Auth/Tenant middleware)
		tenantID := c.GetString(ContextTenantID)
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

		// 3. Final validation after injection
		if err := binding.Validator.ValidateStruct(&input); err != nil {
			_ = c.Error(err).SetType(gin.ErrorTypeBind)
			c.Abort()
			return
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
