package controller

import (
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ClassroomJoinController struct {
	joinService service.ClassroomJoinService
}

func NewClassroomJoinController(joinService service.ClassroomJoinService) *ClassroomJoinController {
	return &ClassroomJoinController{
		joinService: joinService,
	}
}

// PreviewClassroom previews classroom info before joining
// GET /api/classrooms/preview?code=XXX
func (c *ClassroomJoinController) PreviewClassroom(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Invitation code is required", "MISSING_CODE")
		return
	}

	preview, err := c.joinService.PreviewClassroom(code)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Classroom preview retrieved", preview)
}

// JoinClassroom handles join classroom request
// POST /api/classrooms/join
func (c *ClassroomJoinController) JoinClassroom(ctx *gin.Context) {
	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Parse request
	var req dto.JoinClassroomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	// Call service
	result, err := c.joinService.JoinClassroom(user.ID, req)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, result.Message, result)
}

// GetPendingRequests gets pending join requests for a classroom (lecturer only)
// GET /api/classrooms/:id/join-requests
func (c *ClassroomJoinController) GetPendingRequests(ctx *gin.Context) {
	classroomID := ctx.Param("id")

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
	requests, total, err := c.joinService.GetPendingRequests(classroomID, user.ID, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	data := gin.H{
		"requests":  requests,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	}

	dto.SendSuccess(ctx, http.StatusOK, "Pending requests retrieved", data)
}

// ApproveRequest approves a join request
// POST /api/classrooms/join-requests/:id/approve
func (c *ClassroomJoinController) ApproveRequest(ctx *gin.Context) {
	requestID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	err := c.joinService.ApproveRequest(requestID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Join request approved successfully", nil)
}

// RejectRequest rejects a join request
// POST /api/classrooms/join-requests/:id/reject
func (c *ClassroomJoinController) RejectRequest(ctx *gin.Context) {
	requestID := ctx.Param("id")

	// Get authenticated user
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	// Call service
	err := c.joinService.RejectRequest(requestID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Join request rejected successfully", nil)
}
