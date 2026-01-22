package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterClassroomRoutes(router *gin.RouterGroup, controller *controller.ClassroomController, invitationController *controller.ClassroomInvitationController) {
	classrooms := router.Group("/classrooms")
	classrooms.Use(middleware.RequireAuth())
	{
		// CRUD operations
		classrooms.POST("", controller.CreateClassroom)
		classrooms.GET("/my", controller.GetMyClassrooms)
		classrooms.GET("/joined", controller.GetJoinedClassrooms)
		classrooms.GET("/:id", controller.GetClassroomByID)
		classrooms.PUT("/:id", controller.UpdateClassroom)
		classrooms.DELETE("/:id", controller.DeleteClassroom)

		// Get by channel ID
		classrooms.GET("/channel/:id", controller.GetClassroomByChannelID)
		// Status management
		classrooms.PATCH("/:id/status", controller.UpdateClassroomStatus)

		// Leave classroom
		classrooms.POST("/:id/leave", controller.LeaveClassroom)

		// Student management
		classrooms.DELETE("/:id/students/:student_id", controller.RemoveStudentFromClassroom)

		// Co-lecturer management
		classrooms.DELETE("/:id/co-lecturers/:co_lecturer_id", controller.RemoveCoLecturerFromClassroom)

		// Whitelist management
		classrooms.GET("/whitelist-template", controller.DownloadWhitelistTemplate)
		classrooms.GET("/:id/whitelist-student-code", controller.GetWhitelistStudentCode)
		classrooms.POST("/:id/whitelist-student-code", controller.UploadWhitelistStudentCodeJSON)
		classrooms.POST("/:id/whitelist-student-code/upload", controller.UploadWhitelistStudentCodeExcel)
		classrooms.PATCH("/:id/whitelist-student-code", controller.UpdateWhitelistStudentCode)
		classrooms.DELETE("/:id/whitelist-student-code", controller.ClearWhitelistStudentCode)

		// Invitation code
		classrooms.POST("/:id/regenerate-code", controller.RegenerateInvitationCode)

		// Classroom invitations (invite co-lecturers)
		classrooms.POST("/:id/invitations", invitationController.InviteToClassroom)
		classrooms.GET("/:id/invitations", invitationController.GetClassroomInvitations)
	}

	// User's classroom invitations
	invitations := router.Group("/classroom-invitations")
	invitations.Use(middleware.RequireAuth())
	{
		invitations.GET("/my", invitationController.GetMyInvitations)
		invitations.POST("/:id/accept", invitationController.AcceptInvitation)
		invitations.POST("/:id/reject", invitationController.RejectInvitation)
		invitations.DELETE("/:id", invitationController.CancelInvitation)
	}
}
