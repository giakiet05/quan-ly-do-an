package controller

import (
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
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

// InviteToClassroom invites a user to classroom by email
// POST /api/classrooms/:id/invitations
func (c *ClassroomInvitationController) InviteToClassroom(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var req dto.InviteToClassroomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	invitation, err := c.invitationService.InviteByEmail(classroomID, user.ID, req.Email)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Invitation sent successfully", invitation)
}

// AcceptInvitation accepts a classroom invitation
// POST /api/classroom-invitations/:id/accept
func (c *ClassroomInvitationController) AcceptInvitation(ctx *gin.Context) {
	invitationID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	err := c.invitationService.AcceptInvitation(invitationID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation accepted successfully", nil)
}

// RejectInvitation rejects a classroom invitation
// POST /api/classroom-invitations/:id/reject
func (c *ClassroomInvitationController) RejectInvitation(ctx *gin.Context) {
	invitationID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	err := c.invitationService.RejectInvitation(invitationID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation rejected successfully", nil)
}

// CancelInvitation cancels a pending invitation (lecturer/co-lecturer only)
// DELETE /api/classroom-invitations/:id
func (c *ClassroomInvitationController) CancelInvitation(ctx *gin.Context) {
	invitationID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	err := c.invitationService.CancelInvitation(invitationID, user.ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation cancelled successfully", nil)
}

// GetMyInvitations gets pending invitations for the current user
// GET /api/classroom-invitations/my
func (c *ClassroomInvitationController) GetMyInvitations(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var query struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	invitations, total, err := c.invitationService.GetPendingInvitationsForUser(user.ID, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitations retrieved", gin.H{
		"invitations": invitations,
		"total":       total,
		"page":        query.Page,
		"page_size":   query.PageSize,
	})
}

// GetClassroomInvitations gets pending invitations for a classroom
// GET /api/classrooms/:id/invitations
func (c *ClassroomInvitationController) GetClassroomInvitations(ctx *gin.Context) {
	classroomID := ctx.Param("id")

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	user := authUser.(auth.AuthUser)

	var query struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	invitations, total, err := c.invitationService.GetPendingInvitationsForClassroom(classroomID, user.ID, query.Page, query.PageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitations retrieved", gin.H{
		"invitations": invitations,
		"total":       total,
		"page":        query.Page,
		"page_size":   query.PageSize,
	})
}
