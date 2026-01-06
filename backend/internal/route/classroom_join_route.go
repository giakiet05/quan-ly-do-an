package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterClassroomJoinRoutes(router *gin.RouterGroup, controller *controller.ClassroomJoinController) {
	// Public routes
	router.GET("/classrooms/preview", controller.PreviewClassroom)

	// Protected routes (require authentication)
	protected := router.Group("")
	protected.Use(middleware.RequireAuth())
	{
		// Join classroom
		protected.POST("/classrooms/join", controller.JoinClassroom)

		// Manage join requests (lecturer only)
		protected.GET("/classrooms/:id/join-requests", controller.GetPendingRequests)
		protected.POST("/classrooms/join-requests/:id/approve", controller.ApproveRequest)
		protected.POST("/classrooms/join-requests/:id/reject", controller.RejectRequest)
	}
}
