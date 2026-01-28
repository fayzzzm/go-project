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
			// Propagate UserID to Request Context using domain helper
			c.Request = c.Request.WithContext(domain.ContextWithUser(c.Request.Context(), sub))
		}
		if role, ok := claims["role"].(string); ok {
			c.Set(ContextUserRole, role)
		}
		if tenantID, ok := claims["tenant_id"].(string); ok {
			c.Set(ContextTenantID, tenantID)
			// Also put it in the Go context for usecase/repository layers
			c.Request = c.Request.WithContext(domain.ContextWithTenant(c.Request.Context(), tenantID))
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

// GetTenantID retrieves the TenantID from context (Token)
func GetTenantID(c *gin.Context) string {
	if val, exists := c.Get(ContextTenantID); exists {
		if tid, ok := val.(string); ok && tid != "" {
			return tid
		}
	}
	return ""
}

// GetUserRole retrieves the UserRole from context
func GetUserRole(c *gin.Context) string {
	val, exists := c.Get(ContextUserRole)
	if !exists {
		return ""
	}
	return val.(string)
}

// RequireRole checks if the user has the required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient permissions"})
	}
}

// RequireTenant ensures a valid tenant ID is present in the token.
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)

		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Tenant ID is required in the token"})
			return
		}

		// Ensure Go context is also updated (though AuthMiddleware already does this, safety first)
		c.Request = c.Request.WithContext(domain.ContextWithTenant(c.Request.Context(), tenantID))

		c.Next()
	}
}
