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

// GetBodyWithID retrieves the body and injects the ID from the URL parameters.
func GetBodyWithID[T any](c *gin.Context, paramName string) T {
	input := GetBody[T](c)
	v := reflect.ValueOf(&input).Elem()
	if f := v.FieldByName("ID"); f.IsValid() && f.CanSet() {
		f.SetString(c.Param(paramName))
	}
	return input
}
