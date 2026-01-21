package controller

import (
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	projectService service.ProjectService
}

func NewProjectController(projectService service.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

func (c *ProjectController) CreateProjectRound(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.CreateProjectRoundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	round, err := c.projectService.CreateProjectRound(&req, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Project round created successfully", dto.FromProjectRound(round))
}

func (c *ProjectController) CreateProjectRounds(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.CreateProjectRoundsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	rounds, err := c.projectService.CreateProjectRounds(&req, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Project rounds created successfully", dto.FromProjectRounds(rounds))
}

func (c *ProjectController) GetProjectRoundByID(ctx *gin.Context) {
	classroomID := ctx.Param("classroom_id")
	roundID := ctx.Param("round_id")

	round, err := c.projectService.GetProjectRoundByID(classroomID, roundID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Project round retrieved", dto.FromProjectRound(round))
}

func (c *ProjectController) GetProjectRoundsByClassroomID(ctx *gin.Context) {
	classroomID := ctx.Param("classroom_id")

	rounds, err := c.projectService.GetProjectRoundsByClassroomID(classroomID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Project rounds retrieved", dto.FromProjectRounds(rounds))
}

func (c *ProjectController) UpdateProjectRound(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.UpdateProjectRoundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	round, err := c.projectService.UpdateProjectRound(&req, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Project round updated successfully", dto.FromProjectRound(round))
}

func (c *ProjectController) DeleteProjectRound(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	classroomID := ctx.Param("classroom_id")
	roundID := ctx.Param("round_id")

	err := c.projectService.DeleteProjectRound(classroomID, roundID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Project round deleted successfully", nil)
}

// Project endpoints

func (c *ProjectController) CreateProject(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	project, err := c.projectService.CreateProject(&req, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Project created successfully", dto.FromProject(project))
}

func (c *ProjectController) CreateProjects(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.CreateProjectsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	projects, err := c.projectService.CreateProjects(&req, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Projects created successfully", dto.FromProjects(projects))
}

func (c *ProjectController) GetProjectByID(ctx *gin.Context) {
	classroomID := ctx.Param("classroom_id")
	projectID := ctx.Param("project_id")

	project, err := c.projectService.GetProjectByID(classroomID, projectID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Project retrieved", dto.FromProject(project))
}

func (c *ProjectController) GetProjectsByRoundID(ctx *gin.Context) {
	classroomID := ctx.Param("classroom_id")
	roundID := ctx.Param("round_id")

	projects, err := c.projectService.GetProjectsByRoundID(classroomID, roundID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Projects retrieved", dto.FromProjects(projects))
}

func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	project, err := c.projectService.UpdateProject(&req, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Project updated successfully", dto.FromProject(project))
}

func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	classroomID := ctx.Param("classroom_id")
	projectID := ctx.Param("project_id")

	err := c.projectService.DeleteProject(classroomID, projectID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Project deleted successfully", nil)
}

// Report Period endpoints

func (c *ProjectController) CreateReportPeriod(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	classroomID := ctx.Param("classroom_id")
	roundID := ctx.Param("round_id")

	var req dto.CreateReportPeriodRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	period, err := c.projectService.CreateReportPeriod(&req, classroomID, roundID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Report period created successfully", dto.FromReportPeriod(period))
}

func (c *ProjectController) CreateReportPeriods(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	classroomID := ctx.Param("classroom_id")

	var req dto.CreateReportPeriodsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	periods, err := c.projectService.CreateReportPeriods(&req, classroomID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Report periods created successfully", dto.FromReportPeriods(periods))
}

func (c *ProjectController) GetReportPeriodByID(ctx *gin.Context) {
	classroomID := ctx.Param("classroom_id")
	roundID := ctx.Param("round_id")
	periodID := ctx.Param("period_id")

	period, err := c.projectService.GetReportPeriodByID(classroomID, roundID, periodID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report period retrieved", dto.FromReportPeriod(period))
}

func (c *ProjectController) GetReportPeriodsByRoundID(ctx *gin.Context) {
	classroomID := ctx.Param("classroom_id")
	roundID := ctx.Param("round_id")

	periods, err := c.projectService.GetReportPeriodsByRoundID(classroomID, roundID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report periods retrieved", dto.FromReportPeriods(periods))
}

func (c *ProjectController) UpdateReportPeriod(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.UpdateReportPeriodRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	period, err := c.projectService.UpdateReportPeriod(&req, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report period updated successfully", dto.FromReportPeriod(period))
}

func (c *ProjectController) DeleteReportPeriod(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	classroomID := ctx.Param("classroom_id")
	roundID := ctx.Param("round_id")
	periodID := ctx.Param("period_id")

	err := c.projectService.DeleteReportPeriod(classroomID, roundID, periodID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report period deleted successfully", nil)
}
