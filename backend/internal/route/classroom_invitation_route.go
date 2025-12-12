package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupClassroomInvitationRoutes(router *gin.Engine, invitationController *controller.ClassroomInvitationController) {
	// Public routes
	public := router.Group("/api")
	{
		// Download invite template (no auth required)
		public.GET("/classrooms/invite-template", invitationController.DownloadInviteTemplate)
	}

	// Protected routes (require authentication)
	protected := router.Group("/api")
	protected.Use(middleware.RequireAuth())
	{
		// User's pending invitations
		protected.GET("/invitations/my-pending", invitationController.GetMyPendingInvitations)

		// Accept/Reject invitations
		protected.POST("/invitations/:id/accept", invitationController.AcceptInvitation)
		protected.POST("/invitations/:id/reject", invitationController.RejectInvitation)

		// Classroom invitation management (for classroom owner/lecturer)
		protected.POST("/classrooms/:id/invite", invitationController.InviteMembers)
		protected.POST("/classrooms/:id/invite/excel", invitationController.InviteFromExcel)
		protected.GET("/classrooms/:id/invitations", invitationController.GetClassroomInvitations)

		// Cancel invitation (for classroom owner)
		protected.DELETE("/invitations/:id", invitationController.CancelInvitation)
	}
}
