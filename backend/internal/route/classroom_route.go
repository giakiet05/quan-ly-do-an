package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterClassroomRoutes(router *gin.RouterGroup, controller *controller.ClassroomController) {
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
		
		// Status management
		classrooms.PATCH("/:id/status", controller.UpdateClassroomStatus)
		
		// Student management
		classrooms.DELETE("/:id/students/:student_id", controller.RemoveStudentFromClassroom)

		// Whitelist management
		classrooms.GET("/whitelist-template", controller.DownloadWhitelistTemplate)
		classrooms.GET("/:id/whitelist-student-code", controller.GetWhitelistStudentCode)
		classrooms.POST("/:id/whitelist-student-code", controller.UploadWhitelistStudentCodeJSON)
		classrooms.POST("/:id/whitelist-student-code/upload", controller.UploadWhitelistStudentCodeExcel)
		classrooms.PATCH("/:id/whitelist-student-code", controller.UpdateWhitelistStudentCode)
		classrooms.DELETE("/:id/whitelist-student-code", controller.ClearWhitelistStudentCode)
		
		// Invitation code
		classrooms.POST("/:id/regenerate-code", controller.RegenerateInvitationCode)
	}
}
