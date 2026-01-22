package controller

import (
	"fmt"
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

// Project Excel endpoints

// DownloadProjectTemplate generates and downloads an Excel template for project upload
// GET /api/projects/template
func (c *ProjectController) DownloadProjectTemplate(ctx *gin.Context) {
	// Create a new Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Get the default sheet name
	sheetName := "Sheet1"

	// Set headers (Vietnamese)
	f.SetCellValue(sheetName, "A1", "Tên đề tài")
	f.SetCellValue(sheetName, "B1", "Mô tả")
	f.SetCellValue(sheetName, "C1", "Số nhóm")
	f.SetCellValue(sheetName, "D1", "Số SV tối thiểu")
	f.SetCellValue(sheetName, "E1", "Số SV tối đa")

	// Add example rows
	examples := [][]interface{}{
		{"Hệ thống quản lý đồ án", "Xây dựng web app để quản lý đồ án sinh viên", 3, 2, 4},
		{"Ứng dụng di động", "Phát triển app mobile cho sinh viên", 2, 3, 5},
		{"Website thương mại điện tử", "Xây dựng trang web bán hàng online", 1, 2, 3},
	}

	for i, example := range examples {
		rowNum := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), example[0])
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), example[1])
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), example[2])
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), example[3])
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), example[4])
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 30)
	f.SetColWidth(sheetName, "B", "B", 50)
	f.SetColWidth(sheetName, "C", "C", 15)
	f.SetColWidth(sheetName, "D", "D", 20)
	f.SetColWidth(sheetName, "E", "E", 20)

	// Set header style (bold)
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	f.SetCellStyle(sheetName, "A1", "E1", style)

	// Set response headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=projects_template.xlsx")

	// Write to response
	if err := f.Write(ctx.Writer); err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to generate template", "TEMPLATE_ERROR")
		return
	}
}

// UploadProjectsExcel handles bulk project creation from Excel file
// POST /api/projects/classrooms/:classroom_id/rounds/:round_id/upload-excel
func (c *ProjectController) UploadProjectsExcel(ctx *gin.Context) {
	classroomID := ctx.Param("classroom_id")
	roundID := ctx.Param("round_id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Get uploaded file
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "No file uploaded", "NO_FILE")
		return
	}

	// Validate file
	if err := util.ValidateExcelFile(fileHeader); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, err.Error(), "INVALID_FILE")
		return
	}

	// Parse Excel file and extract project data
	projects, parseErrors, err := util.ParseProjectsFromExcel(fileHeader)
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, err.Error(), "PARSE_ERROR")
		return
	}

	// Call service to create projects from Excel data
	createdProjects, err := c.projectService.CreateProjectsFromExcel(classroomID, roundID, user.ID, projects)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	// Prepare response with created projects and any parse errors
	responseData := gin.H{
		"total_projects": len(createdProjects),
		"projects":       dto.FromProjects(createdProjects),
	}

	// Include parse errors if any
	if len(parseErrors) > 0 {
		errorDetails := make([]gin.H, 0)
		for _, parseErr := range parseErrors {
			errorDetails = append(errorDetails, gin.H{
				"row":    parseErr.Row,
				"errors": parseErr.Errors,
			})
		}
		responseData["parse_errors"] = errorDetails
		responseData["total_errors"] = len(parseErrors)
	}

	message := fmt.Sprintf("Successfully created %d projects from Excel", len(createdProjects))
	if len(parseErrors) > 0 {
		message = fmt.Sprintf("Successfully created %d projects. %d rows had errors and were skipped", len(createdProjects), len(parseErrors))
	}

	dto.SendSuccess(ctx, http.StatusCreated, message, responseData)
}
