package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ContextUser represents the authenticated user stored in the Gin context.
type ContextUser struct {
	ID       string
	Username string
	Role     string
}

const contextUserKey = "auth.user"

// AuthMiddleware validates JWT tokens and injects the authenticated user into the context.
func AuthMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be in format 'Bearer <token>'"})
			return
		}

		claims, err := jwtService.ParseToken(strings.TrimSpace(parts[1]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(contextUserKey, ContextUser{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		})
		c.Next()
	}
}

// GetContextUser extracts the authenticated user from the Gin context.
func GetContextUser(c *gin.Context) (ContextUser, bool) {
	value, ok := c.Get(contextUserKey)
	if !ok {
		return ContextUser{}, false
	}
	user, ok := value.(ContextUser)
	return user, ok
}

// RequireRoles ensures the authenticated user has at least one of the provided roles.
func RequireRoles(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		roleSet[strings.ToLower(strings.TrimSpace(role))] = struct{}{}
	}

	return func(c *gin.Context) {
		user, ok := GetContextUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		if _, allowed := roleSet[strings.ToLower(user.Role)]; !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
