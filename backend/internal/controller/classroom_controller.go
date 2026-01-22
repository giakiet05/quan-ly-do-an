package controller

import (
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ClassroomController struct {
	classroomService service.ClassroomService
}

func NewClassroomController(classroomService service.ClassroomService) *ClassroomController {
	return &ClassroomController{
		classroomService: classroomService,
	}
}

// CreateClassroom creates a new classroom
// POST /api/classrooms
func (c *ClassroomController) CreateClassroom(ctx *gin.Context) {
	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse request
	var req dto.CreateClassroomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Call service
	response, err := c.classroomService.CreateClassroom(req, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Classroom created successfully", response)
}

// GetClassroomByID gets classroom by ID
// GET /api/classrooms/:id
func (c *ClassroomController) GetClassroomByID(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	response, err := c.classroomService.GetClassroomByID(classroomID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Classroom retrieved", response)
}

// GET /api/classrooms/channel/:id
func (c *ClassroomController) GetClassroomByChannelID(ctx *gin.Context) {
	channelID := ctx.Param("id")

	response, err := c.classroomService.GetClassroomByChannelID(channelID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}
	dto.SendSuccess(ctx, http.StatusOK, "Classroom retrieved", response)
}

// GetMyClassrooms gets classrooms where user is lecturer
// GET /api/classrooms/my
func (c *ClassroomController) GetMyClassrooms(ctx *gin.Context) {
	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse query params
	var query struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	// Call service
	responses, total, err := c.classroomService.GetClassroomsByLecturer(user.ID, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := gin.H{
		"classrooms": responses,
		"total":      total,
		"page":       query.Page,
		"page_size":  query.PageSize,
	}

	dto.SendSuccess(ctx, http.StatusOK, "Classrooms retrieved", data)
}

// GetJoinedClassrooms gets classrooms where user is a student
// GET /api/classrooms/joined
func (c *ClassroomController) GetJoinedClassrooms(ctx *gin.Context) {
	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse query params
	var query struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Set defaults
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	// Call service
	responses, total, err := c.classroomService.GetClassroomsByStudent(user.ID, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := gin.H{
		"classrooms": responses,
		"total":      total,
		"page":       query.Page,
		"page_size":  query.PageSize,
	}

	dto.SendSuccess(ctx, http.StatusOK, "Joined classrooms retrieved", data)
}

// UpdateClassroom updates classroom information
// PUT /api/classrooms/:id
func (c *ClassroomController) UpdateClassroom(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse request
	var req dto.UpdateClassroomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Call service
	err := c.classroomService.UpdateClassroom(classroomID, user.ID, req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Classroom updated successfully", nil)
}

// UploadWhitelistStudentCodeJSON uploads student code whitelist via JSON
// POST /api/classrooms/:id/whitelist-student-code
func (c *ClassroomController) UploadWhitelistStudentCodeJSON(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse request
	var req dto.UploadWhitelistStudentCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Call service
	err := c.classroomService.UploadWhitelistStudentCode(classroomID, user.ID, req.StudentCodes)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Whitelist updated successfully", nil)
}

// UploadWhitelistStudentCodeExcel uploads student code whitelist via Excel file
// POST /api/classrooms/:id/whitelist-student-code/upload
func (c *ClassroomController) UploadWhitelistStudentCodeExcel(ctx *gin.Context) {
	classroomID := ctx.Param("id")

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

	// Parse Excel file and extract student codes
	studentCodes, err := util.ParseStudentCodesFromExcel(fileHeader)
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, err.Error(), "PARSE_ERROR")
		return
	}

	// Call service to upload whitelist
	err = c.classroomService.UploadWhitelistStudentCode(classroomID, user.ID, studentCodes)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Whitelist uploaded successfully from Excel", gin.H{
		"total_codes": len(studentCodes),
	})
}

// UpdateWhitelistStudentCode adds or removes student codes from whitelist
// PATCH /api/classrooms/:id/whitelist-student-code
func (c *ClassroomController) UpdateWhitelistStudentCode(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse request
	var req dto.UpdateWhitelistStudentCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Call service
	err := c.classroomService.UpdateWhitelistStudentCode(classroomID, user.ID, req.AddCodes, req.RemoveCodes)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Whitelist updated successfully", nil)
}

// ClearWhitelistStudentCode clears all student codes from whitelist
// DELETE /api/classrooms/:id/whitelist-student-code
func (c *ClassroomController) ClearWhitelistStudentCode(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	err := c.classroomService.ClearWhitelistStudentCode(classroomID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Whitelist cleared successfully", nil)
}

// GetWhitelistStudentCode gets current whitelist
// GET /api/classrooms/:id/whitelist-student-code
func (c *ClassroomController) GetWhitelistStudentCode(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	entries, err := c.classroomService.GetWhitelistStudentCode(classroomID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	// Convert to response format
	whitelist := make([]dto.WhitelistEntryResponse, 0, len(entries))
	for _, entry := range entries {
		var joinedBy *string
		var joinedAt *string
		if entry.JoinedBy != nil {
			hex := entry.JoinedBy.Hex()
			joinedBy = &hex
		}
		if entry.JoinedAt != nil {
			formatted := entry.JoinedAt.Format("2006-01-02T15:04:05Z07:00")
			joinedAt = &formatted
		}
		whitelist = append(whitelist, dto.WhitelistEntryResponse{
			StudentCode: entry.StudentCode,
			JoinedBy:    joinedBy,
			JoinedAt:    joinedAt,
		})
	}

	dto.SendSuccess(ctx, http.StatusOK, "Whitelist retrieved", gin.H{
		"whitelist": whitelist,
	})
}

// RegenerateInvitationCode generates a new invitation code
// POST /api/classrooms/:id/regenerate-code
func (c *ClassroomController) RegenerateInvitationCode(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	newCode, err := c.classroomService.RegenerateInvitationCode(classroomID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation code regenerated", dto.RegenerateCodeResponse{
		NewCode: newCode,
	})
}

// UpdateClassroomStatus updates classroom status
// PATCH /api/classrooms/:id/status
func (c *ClassroomController) UpdateClassroomStatus(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse request
	var req dto.UpdateClassroomStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Call service
	err := c.classroomService.UpdateClassroomStatus(classroomID, user.ID, req.Status)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Classroom status updated successfully", nil)
}

// DeleteClassroom deletes a classroom
// DELETE /api/classrooms/:id
func (c *ClassroomController) DeleteClassroom(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	err := c.classroomService.DeleteClassroom(classroomID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Classroom deleted successfully", gin.H{
		"id": classroomID,
	})
}

// RemoveStudentFromClassroom removes a student from classroom
// DELETE /api/classrooms/:id/students/:student_id
func (c *ClassroomController) RemoveStudentFromClassroom(ctx *gin.Context) {
	classroomID := ctx.Param("id")
	studentID := ctx.Param("student_id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	err := c.classroomService.RemoveStudentFromClassroom(classroomID, user.ID, studentID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Student removed from classroom successfully", gin.H{
		"student_id": studentID,
	})
}

// RemoveCoLecturerFromClassroom removes a co-lecturer from classroom (lecturer only)
// DELETE /api/classrooms/:id/co-lecturers/:co_lecturer_id
func (c *ClassroomController) RemoveCoLecturerFromClassroom(ctx *gin.Context) {
	classroomID := ctx.Param("id")
	coLecturerID := ctx.Param("co_lecturer_id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	err := c.classroomService.RemoveCoLecturerFromClassroom(classroomID, user.ID, coLecturerID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Co-lecturer removed from classroom successfully", gin.H{
		"co_lecturer_id": coLecturerID,
	})
}

// LeaveClassroom allows a student or co-lecturer to leave the classroom
// POST /api/classrooms/:id/leave
func (c *ClassroomController) LeaveClassroom(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	err := c.classroomService.LeaveClassroom(classroomID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Left classroom successfully", nil)
}

// DownloadWhitelistTemplate generates and downloads an Excel template for whitelist upload
// GET /api/classrooms/whitelist-template
func (c *ClassroomController) DownloadWhitelistTemplate(ctx *gin.Context) {
	// Create a new Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Get the default sheet name
	sheetName := "Sheet1"

	// Set header
	f.SetCellValue(sheetName, "A1", "Student Code")

	// Add example rows
	examples := []string{"2211234", "2211235", "2211236"}
	for i, code := range examples {
		cellAddr, _ := excelize.CoordinatesToCellName(1, i+2)
		f.SetCellValue(sheetName, cellAddr, code)
	}

	// Set column width
	f.SetColWidth(sheetName, "A", "A", 20)

	// Set header style (bold)
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	f.SetCellStyle(sheetName, "A1", "A1", style)

	// Set response headers for file download
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=whitelist_template.xlsx")

	// Write to response
	if err := f.Write(ctx.Writer); err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to generate template", "TEMPLATE_ERROR")
		return
	}
}
