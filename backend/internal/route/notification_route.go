package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(rg *gin.RouterGroup, c *controller.NotificationController) {
	notifications := rg.Group("/notifications")
	notifications.Use(middleware.RequireAuth()) // All notification routes require authentication
	{
		notifications.GET("", c.GetNotifications)
		notifications.PUT("/read-all", c.MarkAllAsRead)
	}
}
