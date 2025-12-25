package controller

import (
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/excel"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ClassroomInvitationController struct {
	invitationService service.ClassroomInvitationService
}

func NewClassroomInvitationController(invitationService service.ClassroomInvitationService) *ClassroomInvitationController {
	return &ClassroomInvitationController{
		invitationService: invitationService,
	}
}

// InviteMembers invites multiple emails to a classroom (manual input)
// POST /api/classrooms/:id/invite
func (c *ClassroomInvitationController) InviteMembers(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse request
	var req dto.InviteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Call service
	result := c.invitationService.InviteMembers(classroomID, user.ID, req.Emails)

	dto.SendSuccess(ctx, http.StatusOK, "Invitations processed", result)
}

// InviteFromExcel invites emails from Excel file
// POST /api/classrooms/:id/invite/excel
func (c *ClassroomInvitationController) InviteFromExcel(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Get uploaded file
	file, err := ctx.FormFile("file")
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, "Missing file", "MISSING_FILE")
		return
	}

	// Validate file type
	if !isExcelFile(file.Filename) {
		dto.SendError(ctx, http.StatusBadRequest, "Only Excel files (.xlsx, .xls) are allowed", "INVALID_FILE_TYPE")
		return
	}

	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024)
	if file.Size > maxSize {
		dto.SendError(ctx, http.StatusBadRequest, "File size exceeds 5MB", "FILE_TOO_LARGE")
		return
	}

	// Parse Excel file
	emails, parseErrors, err := excel.ParseEmailsFromExcel(file)
	if err != nil {
		dto.SendError(ctx, http.StatusBadRequest, err.Error(), "PARSE_ERROR")
		return
	}

	// Call service to invite members
	result := c.invitationService.InviteMembers(classroomID, user.ID, emails)

	// Merge parse errors into failed results
	for _, parseErr := range parseErrors {
		result.Failed = append(result.Failed, dto.InviteError{
			Row:   parseErr.Row,
			Email: parseErr.Value,
			Error: parseErr.Error,
		})
	}

	dto.SendSuccess(ctx, http.StatusOK, "Excel file processed", result)
}

// GetMyPendingInvitations gets all pending invitations for current user
// GET /api/invitations/my-pending
func (c *ClassroomInvitationController) GetMyPendingInvitations(ctx *gin.Context) {
	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	invitations, err := c.invitationService.GetMyPendingInvitations(user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Pending invitations retrieved", invitations)
}

// GetClassroomInvitations gets invitations for a classroom (for owner/admin)
// GET /api/classrooms/:id/invitations
func (c *ClassroomInvitationController) GetClassroomInvitations(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	// Parse query params
	var query struct {
		Status   *string `form:"status"`
		Page     int     `form:"page"`
		PageSize int     `form:"pageSize"`
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

	// Parse status filter
	var statusFilter *model.InvitationStatus
	if query.Status != nil && *query.Status != "" {
		status := model.InvitationStatus(*query.Status)
		statusFilter = &status
	}

	// Call service
	invitations, total, err := c.invitationService.GetClassroomInvitations(classroomID, statusFilter, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := gin.H{
		"invitations": invitations,
		"total":       total,
		"page":        query.Page,
		"page_size":   query.PageSize,
	}

	dto.SendSuccess(ctx, http.StatusOK, "Classroom invitations retrieved", data)
}

// AcceptInvitation accepts an invitation
// POST /api/invitations/:id/accept
func (c *ClassroomInvitationController) AcceptInvitation(ctx *gin.Context) {
	invitationID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	result, err := c.invitationService.AcceptInvitation(invitationID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation accepted successfully", result)
}

// RejectInvitation rejects an invitation
// POST /api/invitations/:id/reject
func (c *ClassroomInvitationController) RejectInvitation(ctx *gin.Context) {
	invitationID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	err := c.invitationService.RejectInvitation(invitationID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation rejected successfully", nil)
}

// CancelInvitation cancels an invitation (by classroom owner)
// DELETE /api/invitations/:id
func (c *ClassroomInvitationController) CancelInvitation(ctx *gin.Context) {
	invitationID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	err := c.invitationService.CancelInvitation(invitationID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation canceled successfully", nil)
}

// DownloadInviteTemplate downloads Excel template for classroom invitations
// GET /api/classrooms/invite-template
func (c *ClassroomInvitationController) DownloadInviteTemplate(ctx *gin.Context) {
	// Create template
	f, err := excel.CreateInviteTemplate()
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to create template", "TEMPLATE_ERROR")
		return
	}
	defer f.Close()

	// Set response headers
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=classroom_invite_template.xlsx")

	// Write to response
	if err := f.Write(ctx.Writer); err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to write template", "WRITE_ERROR")
		return
	}
}

// Helper functions

func isExcelFile(filename string) bool {
	return len(filename) > 5 && (filename[len(filename)-5:] == ".xlsx" || filename[len(filename)-4:] == ".xls")
}
