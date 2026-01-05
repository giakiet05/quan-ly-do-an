package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterGroupRoutes(rg *gin.RouterGroup, c *controller.GroupController) {
	groups := rg.Group("/groups")
	// Protected routes (require authentication)
	groups.Use(middleware.RequireAuth())
	{
		// Group CRUD operations
		groups.POST("", c.CreateGroup)
		groups.GET("/:group_id", c.GetGroupByID)
		groups.GET("", c.GetGroupsFilter)
		groups.PUT("", c.UpdateGroup)
		groups.DELETE("/:group_id", c.DeleteGroup)

		groups.PUT("leave/:group_id", c.LeaveGroup)

		groups.POST("join", c.CreateJoinGroupRequest)
		groups.PUT("join", c.UpdateJoinGroupRequest)

		groups.PUT("invitation", c.UpdateGroupInvitation)
		groups.POST("invitation", c.CreateGroupInvitation)

		// Task operations
		groups.POST("/tasks", c.CreateTask)
		groups.PUT("/tasks", c.UpdateTask)
		groups.DELETE("/:group_id/tasks/:task_id", c.DeleteTask)

		// Report operations
		groups.POST("/reports", c.CreateReport)
		groups.PUT("/reports", c.UpdateReport)
		groups.DELETE("/:group_id/reports/:report_id", c.DeleteReport)

		// Report feedback operations
		groups.POST("/reports/feedback", c.CreateReportFeedback)
		groups.PUT("/:group_id/reports/:report_id/feedback", c.UpdateReportFeedback)
		groups.DELETE("/:group_id/reports/:report_id/feedback", c.DeleteReportFeedback)
	}
}
