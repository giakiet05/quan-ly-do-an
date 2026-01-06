package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterClassPostRoutes(router *gin.RouterGroup, controller *controller.ClassPostController) {
	// Posts within a classroom
	router.POST("/classrooms/:id/posts", middleware.RequireAuth(), controller.CreatePost)
	router.GET("/classrooms/:id/posts", middleware.RequireAuth(), controller.GetPosts)

	// Individual post operations
	posts := router.Group("/classrooms/posts")
	posts.Use(middleware.RequireAuth())
	{
		// File operations
		posts.POST("/upload-attachments", controller.UploadAttachments)
		posts.DELETE("/attachments/:public_id", controller.DeleteAttachment)
		
		// Post CRUD
		posts.GET("/:post_id", controller.GetPost)
		posts.PUT("/:post_id", controller.UpdatePost)
		posts.DELETE("/:post_id", controller.DeletePost)
		posts.PATCH("/:post_id/pin", controller.TogglePinPost)
		
		// Attachment management
		posts.POST("/:post_id/attachments", controller.AddAttachment)
		posts.DELETE("/:post_id/attachments", controller.RemoveAttachment)
	}
}
