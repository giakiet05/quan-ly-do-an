package route

import (
	"github.com/giakiet05/quan-ly-do-an/backend/internal/controller"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(rg *gin.RouterGroup, c *controller.ProjectController) {
	projects := rg.Group("/projects")

	// Protected routes (require authentication)
	projects.Use(middleware.RequireAuth())
	{
		// ===== Report Period operations =====
		projects.POST("/classrooms/:classroom_id/rounds/:round_id/report-periods", c.CreateReportPeriod)
		projects.POST("/classrooms/:classroom_id/report-periods/bulk", c.CreateReportPeriods)

		projects.GET("/classrooms/:classroom_id/rounds/:round_id/report-periods/:period_id", c.GetReportPeriodByID)
		projects.GET("/classrooms/:classroom_id/rounds/:round_id/report-periods", c.GetReportPeriodsByRoundID)

		projects.PUT("/report-periods", c.UpdateReportPeriod)
		projects.DELETE("/classrooms/:classroom_id/rounds/:round_id/report-periods/:period_id", c.DeleteReportPeriod)

		// ===== Project operations =====
		projects.POST("", c.CreateProject)
		projects.POST("/bulk", c.CreateProjects)

		projects.GET("/classrooms/:classroom_id/:project_id", c.GetProjectByID)
		projects.GET("/classrooms/:classroom_id/rounds/:round_id/projects", c.GetProjectsByRoundID)

		projects.PUT("", c.UpdateProject)
		projects.DELETE("/classrooms/:classroom_id/:project_id", c.DeleteProject)

		// ===== Project Round operations =====
		projects.POST("/rounds", c.CreateProjectRound)
		projects.POST("/rounds/bulk", c.CreateProjectRounds)

		projects.GET("/classrooms/:classroom_id/rounds/:round_id", c.GetProjectRoundByID)
		projects.GET("/classrooms/:classroom_id/rounds", c.GetProjectRoundsByClassroomID)

		projects.PUT("/rounds", c.UpdateProjectRound)
		projects.DELETE("/classrooms/:classroom_id/rounds/:round_id", c.DeleteProjectRound)
	}
}
