package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c *controller.UserController) {
	users := rg.Group("/users")

	// Public routes
	users.GET("/", c.GetUsers)

	// Routes for the currently authenticated user ("me")
	me := users.Group("/me")
	me.Use(middleware.RequireAuth())
	{
		me.GET("", c.GetMyProfile)
		me.PUT("/password", c.ChangePassword)
		me.POST("/avatar", c.UploadAvatar)
		me.DELETE("/avatar", c.DeleteAvatar)
	}
}
