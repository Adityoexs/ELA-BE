package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const claimsKey = "auth_claims"

// Middleware returns a Gin handler that enforces JWT authentication.
// On success it stores the validated *Claims under the key "auth_claims" in the
// Gin context so downstream handlers (and future RBAC middleware) can read them.
func Middleware(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must use Bearer scheme"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)
		claims, err := svc.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(claimsKey, claims)
		c.Next()
	}
}

// GetClaims retrieves the validated JWT claims stored by Middleware.
// Returns nil if the claims are not present (e.g. on unprotected routes) or
// if the stored value is not of the expected *Claims type.
func GetClaims(c *gin.Context) *Claims {
	val, exists := c.Get(claimsKey)
	if !exists {
		return nil
	}
	claims, ok := val.(*Claims)
	if !ok {
		return nil
	}
	return claims
}
