package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	HeaderAuthorization = "Authorization"
	ContextUserID       = "user_id"
	ContextUserRole     = "user_role"
	ContextTenantID     = "tenant_id"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return []byte(secret)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(HeaderAuthorization)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Bearer prefix
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		tokenString = parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return getJWTSecret(), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		if sub, ok := claims["sub"].(string); ok {
			c.Set(ContextUserID, sub)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set(ContextUserRole, role)
		}
		if tenantID, ok := claims["tenant_id"].(string); ok {
			c.Set(ContextTenantID, tenantID)
		}

		c.Next()
	}
}

// GetUserID retrieves the UserID from context
func GetUserID(c *gin.Context) (string, error) {
	val, exists := c.Get(ContextUserID)
	if !exists {
		return "", domain.ErrUnauthorized
	}
	return val.(string), nil
}

// GetTenantID retrieves the TenantID from context (Token) or Header
func GetTenantID(c *gin.Context) string {
	// 1. Try Token
	if val, exists := c.Get(ContextTenantID); exists {
		if tid, ok := val.(string); ok && tid != "" {
			return tid
		}
	}
	// 2. Try Header
	return c.GetHeader("X-Tenant-ID")
}
