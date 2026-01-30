package controller

import (
	"encoding/json"
	"net/http"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/auth"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/cloudinary"
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
	var req dto.CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	group, err := g.groupService.CreateGroup(&req, authUser.(auth.AuthUser).ID)
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
	var query dto.GetGroupsFilterQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	groups, err := g.groupService.GetGroupsFilter(&query, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Groups retrieved successfully", dto.FromGroups(groups))
}

func (g *GroupController) UpdateGroup(ctx *gin.Context) {
	var req dto.UpdateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	group, err := g.groupService.UpdateGroup(&req, authUser.(auth.AuthUser).ID)
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
	userID := ctx.Param("user_id")
	if groupID == "" || userID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID and User ID are required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.RemoveMemberFromGroup(groupID, userID, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Left group successfully", nil)
}

// Join Request Operations

func (g *GroupController) CreateJoinRequest(ctx *gin.Context) {
	var req dto.CreateJoinGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	joinRequest, err := g.groupService.CreateJoinRequest(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Join request created successfully", joinRequest)
}

func (g *GroupController) AcceptJoinRequest(ctx *gin.Context) {
	var req dto.UpdateJoinGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.AcceptJoinRequest(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Join request accepted successfully", nil)
}

func (g *GroupController) RejectJoinRequest(ctx *gin.Context) {
	var req dto.UpdateJoinGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.RejectJoinRequest(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Join request rejected successfully", nil)
}

// Invitation Operations

func (g *GroupController) InviteToGroup(ctx *gin.Context) {
	var req dto.CreateGroupInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	invitation, err := g.groupService.InviteToGroup(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Invitation sent successfully", invitation)
}

func (g *GroupController) AcceptInvitation(ctx *gin.Context) {
	var req dto.UpdateGroupInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.AcceptInvitation(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation accepted successfully", nil)
}

func (g *GroupController) RejectInvitation(ctx *gin.Context) {
	var req dto.UpdateGroupInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := g.groupService.RejectInvitation(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Invitation rejected successfully", nil)
}

// Task Operations

func (g *GroupController) CreateTask(ctx *gin.Context) {
	var req dto.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	task, err := g.groupService.CreateTask(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Task created successfully", task)
}

func (g *GroupController) UpdateTask(ctx *gin.Context) {
	var req dto.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	task, err := g.groupService.UpdateTask(&req, authUser.(auth.AuthUser).ID)
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
	// Parse form data
	classroomID := ctx.PostForm("classroom_id")
	groupID := ctx.PostForm("group_id")
	projectRoundID := ctx.PostForm("project_round_id")
	reportPeriodID := ctx.PostForm("report_period_id")
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")

	if classroomID == "" || groupID == "" || projectRoundID == "" || reportPeriodID == "" || title == "" || content == "" {
		dto.SendError(ctx, http.StatusBadRequest, "All fields are required", apperror.ErrBadRequest.Code)
		return
	}

	// Get uploaded files
	form, err := ctx.MultipartForm()
	var attachments []dto.AttachmentUpload
	uploadedPublicIDs := []string{}

	if err == nil && form != nil {
		files := form.File["files"]

		// Upload files to Cloudinary
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				// Rollback: delete uploaded files
				for _, pid := range uploadedPublicIDs {
					_, _ = cloudinary.Delete(pid)
				}
				dto.SendError(ctx, http.StatusInternalServerError, "Failed to open file", "FILE_OPEN_FAILED")
				return
			}

			result, err := cloudinary.Upload(file)
			file.Close()

			if err != nil {
				// Rollback: delete uploaded files
				for _, pid := range uploadedPublicIDs {
					_, _ = cloudinary.Delete(pid)
				}
				dto.SendError(ctx, http.StatusInternalServerError, "Failed to upload file", "FILE_UPLOAD_FAILED")
				return
			}

			uploadedPublicIDs = append(uploadedPublicIDs, result.PublicID)
			attachments = append(attachments, dto.AttachmentUpload{
				FileName: fileHeader.Filename,
				FileURL:  result.SecureURL,
				PublicID: result.PublicID,
				FileSize: fileHeader.Size,
				MimeType: fileHeader.Header.Get("Content-Type"),
			})
		}
	}

	// Create report request
	req := &dto.CreateReportRequest{
		ClassroomID:    classroomID,
		GroupID:        groupID,
		ProjectRoundID: projectRoundID,
		ReportPeriodID: reportPeriodID,
		Title:          title,
		Content:        content,
		Attachments:    attachments,
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		// Rollback: delete uploaded files
		for _, pid := range uploadedPublicIDs {
			_, _ = cloudinary.Delete(pid)
		}
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	report, err := g.groupService.CreateReport(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		// Rollback: delete uploaded files
		for _, pid := range uploadedPublicIDs {
			_, _ = cloudinary.Delete(pid)
		}
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Report created successfully", report)
}

func (g *GroupController) UpdateReport(ctx *gin.Context) {
	// Parse form data
	groupID := ctx.PostForm("group_id")
	reportID := ctx.PostForm("report_id")
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")
	filesToRemoveStr := ctx.PostForm("files_to_remove")

	if groupID == "" || reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID and Report ID are required", apperror.ErrBadRequest.Code)
		return
	}

	// Parse files to remove
	var filesToRemove []string
	if filesToRemoveStr != "" {
		if err := json.Unmarshal([]byte(filesToRemoveStr), &filesToRemove); err != nil {
			dto.SendError(ctx, http.StatusBadRequest, "Invalid files_to_remove format", apperror.ErrBadRequest.Code)
			return
		}
	}

	// Get uploaded files
	form, err := ctx.MultipartForm()
	var attachmentsToAdd []dto.AttachmentUpload
	uploadedPublicIDs := []string{}

	if err == nil && form != nil {
		files := form.File["files"]

		// Upload files to Cloudinary
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				// Rollback: delete uploaded files
				for _, pid := range uploadedPublicIDs {
					_, _ = cloudinary.Delete(pid)
				}
				dto.SendError(ctx, http.StatusInternalServerError, "Failed to open file", "FILE_OPEN_FAILED")
				return
			}

			result, err := cloudinary.Upload(file)
			file.Close()

			if err != nil {
				// Rollback: delete uploaded files
				for _, pid := range uploadedPublicIDs {
					_, _ = cloudinary.Delete(pid)
				}
				dto.SendError(ctx, http.StatusInternalServerError, "Failed to upload file", "FILE_UPLOAD_FAILED")
				return
			}

			uploadedPublicIDs = append(uploadedPublicIDs, result.PublicID)
			attachmentsToAdd = append(attachmentsToAdd, dto.AttachmentUpload{
				FileName: fileHeader.Filename,
				FileURL:  result.SecureURL,
				PublicID: result.PublicID,
				FileSize: fileHeader.Size,
				MimeType: fileHeader.Header.Get("Content-Type"),
			})
		}
	}

	// Build request
	req := &dto.UpdateReportRequest{
		GroupID:          groupID,
		ReportID:         reportID,
		AttachmentsToAdd: attachmentsToAdd,
		FilesToRemove:    filesToRemove,
	}

	if title != "" {
		req.Title = &title
	}
	if content != "" {
		req.Content = &content
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		// Rollback: delete uploaded files
		for _, pid := range uploadedPublicIDs {
			_, _ = cloudinary.Delete(pid)
		}
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	report, err := g.groupService.UpdateReport(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		// Rollback: delete uploaded files
		for _, pid := range uploadedPublicIDs {
			_, _ = cloudinary.Delete(pid)
		}
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
	var req dto.CreateReportFeedbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	feedback, err := g.groupService.CreateReportFeedback(&req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Report feedback created successfully", feedback)
}

func (g *GroupController) UpdateReportFeedback(ctx *gin.Context) {
	var req dto.UpdateReportFeedbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	groupID := ctx.Param("group_id")
	reportID := ctx.Param("report_id")

	if groupID == "" || reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Group ID and Report ID are required", apperror.ErrBadRequest.Code)
		return
	}

	err := g.groupService.UpdateReportFeedback(&req, groupID, reportID)
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
