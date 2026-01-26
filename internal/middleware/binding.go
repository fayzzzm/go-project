package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// BodyContextKey is the key used to store the parsed body in the context.
const BodyContextKey = "request_body"

// Validatable interface that structs can implement if they have self-validation.
type Validatable interface {
	Validate() (bool, string)
}

// BindJSON is a generic middleware that binds the request body to a struct of type T.
// It also checks if the struct implements Validatable and calls its Validate method.
func BindJSON[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input T
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
			c.Abort()
			return
		}

		// Check if T (or *T) implements Validatable
		if v, ok := any(&input).(Validatable); ok {
			if ok, msg := v.Validate(); !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				c.Abort()
				return
			}
		} else if v, ok := any(input).(Validatable); ok {
			if ok, msg := v.Validate(); !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				c.Abort()
				return
			}
		}

		// Inject X-Tenant-ID if struct has TenantID field and header is present
		if tenantID := c.GetHeader("X-Tenant-ID"); tenantID != "" {
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
