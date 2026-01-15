package middleware

import (
	"net/http"
	"strings"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/gin-gonic/gin"
)

// RequireAuth parse access token từ Authorization header và nhét AuthUser vào context
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			dto.SendError(c, http.StatusUnauthorized, apperror.ErrInvalidToken.Message, apperror.ErrInvalidToken.Code)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			dto.SendError(c, http.StatusUnauthorized, apperror.ErrInvalidToken.Message, apperror.ErrInvalidToken.Code)
			c.Abort()
			return
		}

		token := parts[1]
		user, err := auth.ParseAccessToken(token)
		if err != nil {
			dto.SendError(c, http.StatusUnauthorized, apperror.Message(err), apperror.Code(err))
			c.Abort()
			return
		}

		c.Set("authUser", user)
		c.Next()
	}
}

// RequireAuthSocket parse access token từ query param (cho WebSocket) và nhét AuthUser vào context
func RequireAuthSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			dto.SendError(c, http.StatusUnauthorized, apperror.ErrInvalidToken.Message, apperror.ErrInvalidToken.Code)
			c.Abort()
			return
		}

		user, err := auth.ParseAccessToken(token)
		if err != nil {
			dto.SendError(c, http.StatusUnauthorized, apperror.Message(err), apperror.Code(err))
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
			dto.SendError(c, http.StatusUnauthorized, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
			c.Abort()
			return
		}

		user, ok := val.(auth.AuthUser)
		if !ok {
			dto.SendError(c, http.StatusInternalServerError, apperror.ErrInternal.Message, apperror.ErrInternal.Code)
			c.Abort()
			return
		}

		if user.Role != "admin" {
			dto.SendError(c, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
			c.Abort()
			return
		}

		c.Next()
	}
}
