package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(rg *gin.RouterGroup, c *controller.NotificationController) {
	notifications := rg.Group("/notifications")
	notifications.Use(middleware.RequireAuth())
	{
		notifications.GET("", c.GetNotificationsRecipientID)
		notifications.GET("/classroom/:classroom_id", c.GetNotificationsByClassroomID)
		notifications.POST("/:notification_id/read", c.MarkAsRead)
		notifications.PUT("/read-all", c.MarkAllAsRead)
		notifications.DELETE("/:notification_id", c.DeleteNotification)
	}
}
