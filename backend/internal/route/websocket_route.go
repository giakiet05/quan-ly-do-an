package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterWebSocketRoutes(rg *gin.RouterGroup, c *controller.WebSocketController) {
	ws := rg.Group("/ws")
	ws.Use(middleware.RequireAuthSocket()) // WebSocket connections get token from query param
	{
		ws.GET("", c.HandleConnections)
	}
}
