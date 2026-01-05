package controller

import (
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type GroupController struct {
	groupService service.GroupService
}

func NewGroupController(groupService service.GroupService) *GroupController {
	return &GroupController{
		groupService: groupService,
	}
}

// Group CRUD Operations

func (g *GroupController) CreateGroup(ctx *gin.Context) {
	var req *dto.CreateGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	group, err := g.groupService.CreateGroup(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Group created successfully", dto.FromGroup(group))
}

func (g *GroupController) GetGroupByID(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	if groupID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	group, err := g.groupService.GetGroupByID(groupID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Group retrieved successfully", dto.FromGroup(group))
}

func (g *GroupController) GetGroupsFilter(ctx *gin.Context) {
	var query *dto.GetGroupsFilterQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	groups, err := g.groupService.GetGroupsFilter(query, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Groups retrieved successfully", dto.FromGroups(groups))
}

func (g *GroupController) UpdateGroup(ctx *gin.Context) {
	var req *dto.UpdateGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	group, err := g.groupService.UpdateGroup(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Group updated successfully", dto.FromGroup(group))
}

func (g *GroupController) DeleteGroup(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	if groupID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.DeleteGroup(groupID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Group deleted successfully", nil)
}

func (g *GroupController) LeaveGroup(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	if groupID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.LeaveGroup(groupID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Leave group successfully", nil)
}

func (g *GroupController) CreateJoinGroupRequest(ctx *gin.Context) {
	var req *dto.CreateJoinGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.CreateJoinGroupRequest(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Join group request created successfully", nil)
}

func (g *GroupController) UpdateJoinGroupRequest(ctx *gin.Context) {
	var req *dto.UpdateJoinGroupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.UpdateJoinGroupRequest(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Join group request updated successfully", nil)
}

func (g *GroupController) CreateGroupInvitation(ctx *gin.Context) {
	var req *dto.CreateGroupInvitationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.CreateGroupInvitation([]string{req.RecipientID}, req.GroupID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Group invitation created successfully", nil)
}

func (g *GroupController) UpdateGroupInvitation(ctx *gin.Context) {
	var req *dto.UpdateGroupInvitationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.UpdateGroupInvitation(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Group invitation updated successfully", nil)
}

// Task Operations

func (g *GroupController) CreateTask(ctx *gin.Context) {
	var req *dto.CreateTaskRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	task, err := g.groupService.CreateTask(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Task created successfully", task)
}

func (g *GroupController) UpdateTask(ctx *gin.Context) {
	var req *dto.UpdateTaskRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	task, err := g.groupService.UpdateTask(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Task updated successfully", task)
}

func (g *GroupController) DeleteTask(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	taskID := ctx.Param("task_id")

	if groupID == "" || taskID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID and Task ID are required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.DeleteTask(groupID, taskID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Task deleted successfully", nil)
}

// Report Operations

func (g *GroupController) CreateReport(ctx *gin.Context) {
	var req *dto.CreateReportRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	report, err := g.groupService.CreateReport(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Report created successfully", report)
}

func (g *GroupController) UpdateReport(ctx *gin.Context) {
	var req *dto.UpdateReportRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	report, err := g.groupService.UpdateReport(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report updated successfully", report)
}

func (g *GroupController) DeleteReport(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	reportID := ctx.Param("report_id")

	if groupID == "" || reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID and Report ID are required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.DeleteReport(groupID, reportID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report deleted successfully", nil)
}

// Report Feedback Operations

func (g *GroupController) CreateReportFeedback(ctx *gin.Context) {
	var req *dto.CreateReportFeedbackRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	feedback, err := g.groupService.CreateReportFeedback(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Report feedback created successfully", feedback)
}

func (g *GroupController) UpdateReportFeedback(ctx *gin.Context) {
	var req *dto.UpdateReportFeedbackRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	groupID := ctx.Param("group_id")
	reportID := ctx.Param("report_id")

	if groupID == "" || reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID and Report ID are required", apperror.ErrBadRequest.Code)
		return
	}

	err := g.groupService.UpdateReportFeedback(req, groupID, reportID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report feedback updated successfully", nil)
}

func (g *GroupController) DeleteReportFeedback(ctx *gin.Context) {
	groupID := ctx.Param("group_id")
	reportID := ctx.Param("report_id")

	if groupID == "" || reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID and Report ID are required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.DeleteReportFeedback(groupID, reportID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report feedback deleted successfully", nil)
}
