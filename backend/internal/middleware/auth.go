package middleware

import (
	"net/http"
	"strings"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// RequireAuth parse access token và nhét AuthUser vào context
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		user, err := auth.ParseAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
				"debug": err.Error(), // TEMP: Debug info
			})
			c.Abort()
			return
		} // Load user settings from DB once per request

		// Nhét user vào context with settings cached
		c.Set("authUser", user)
		c.Next()
	}
}

func RequireAuthSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token in query parameter"})
			c.Abort()
			return
		}

		user, err := auth.ParseAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
				"debug": err.Error(),
			})
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		user, ok := val.(auth.AuthUser)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth context"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
