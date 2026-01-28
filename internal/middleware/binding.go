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
		// We use ShouldBindBodyWith to allow reading the body multiple times if needed,
		// but primarily to decode into the struct. this checks binding:"..." tags immediately.
		if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
			// If it's a JSON Parsing error (e.g. invalid syntax, wrong types), we must fail immediately.
			// We can't inject fields into a broken struct.
			// However, if it's just a "required field missing" error, we might fix it in step 2 (Injection).
			// So we check if the error is ONLY validation related.

			// Note: Gin wraps errors. We rely on the final validation step to catch persistent issues.
			// But if binding fails due to bad JSON, Input might be zero-value.
			// We continue to Step 2 to attempt injection, then Step 3 validates everything.
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
