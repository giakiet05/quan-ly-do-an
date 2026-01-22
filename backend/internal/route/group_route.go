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
		groups.POST("", c.CreateGroup)
		groups.GET("/:group_id", c.GetGroupByID)
		groups.GET("", c.GetGroupsFilter)
		groups.PUT("", c.UpdateGroup)
		groups.DELETE("/:group_id", c.DeleteGroup)
		groups.PUT(":group_id/remove/:user_id", c.LeaveGroup)

		groups.POST("/join-requests", c.CreateJoinRequest)
		groups.PUT("/join-requests/accept", c.AcceptJoinRequest)
		groups.PUT("/join-requests/reject", c.RejectJoinRequest)

		groups.POST("/invitations", c.InviteToGroup)
		groups.PUT("/invitations/accept", c.AcceptInvitation)
		groups.PUT("/invitations/reject", c.RejectInvitation)

		groups.POST("/tasks", c.CreateTask)
		groups.PUT("/tasks", c.UpdateTask)
		groups.DELETE("/:group_id/tasks/:task_id", c.DeleteTask)

		groups.POST("/reports", c.CreateReport)
		groups.PUT("/reports", c.UpdateReport)
		groups.DELETE("/:group_id/reports/:report_id", c.DeleteReport)

		groups.POST("/reports/feedback", c.CreateReportFeedback)
		groups.PUT("/:group_id/reports/:report_id/feedback", c.UpdateReportFeedback)
		groups.DELETE("/:group_id/reports/:report_id/feedback", c.DeleteReportFeedback)
	}
}
