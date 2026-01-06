package middleware

import (
	"net/http"
	"strings"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/gin-gonic/gin"
)

// RequireAuth parse access token và nhét AuthUser vào context
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			dto.SendError(c, http.StatusUnauthorized, "Missing Authorization header", "AUTH_MISSING_HEADER")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			dto.SendError(c, http.StatusUnauthorized, "Invalid Authorization header format", "AUTH_INVALID_FORMAT")
			c.Abort()
			return
		}

		token := parts[1]
		user, err := auth.ParseAccessToken(token)
		if err != nil {
			dto.SendError(c, http.StatusUnauthorized, "Invalid or expired token", "AUTH_INVALID_TOKEN")
			c.Abort()
			return
		}

		// Nhét user vào context with settings cached
		c.Set("authUser", user)
		c.Next()
	}
}

func RequireAuthSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			dto.SendError(c, http.StatusUnauthorized, "Missing token in query parameter", "AUTH_MISSING_TOKEN")
			c.Abort()
			return
		}

		user, err := auth.ParseAccessToken(token)
		if err != nil {
			dto.SendError(c, http.StatusUnauthorized, "Invalid or expired token", "AUTH_INVALID_TOKEN")
			c.Abort()
			return
		}

		c.Set("authUser", user)
		c.Next()
	}
}

// RequireAdmin check role admin
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("authUser")
		if !exists {
			dto.SendError(c, http.StatusUnauthorized, "Not authenticated", "AUTH_NOT_AUTHENTICATED")
			c.Abort()
			return
		}

		user, ok := val.(auth.AuthUser)
		if !ok {
			dto.SendError(c, http.StatusInternalServerError, "Invalid auth context", "AUTH_INVALID_CONTEXT")
			c.Abort()
			return
		}

		if user.Role != "admin" {
			dto.SendError(c, http.StatusForbidden, "Admin access required", "AUTH_ADMIN_REQUIRED")
			c.Abort()
			return
		}

		c.Next()
	}
}
