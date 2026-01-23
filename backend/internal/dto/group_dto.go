package dto

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

type CreateGroupRequest struct {
	ClassroomID    string `json:"classroom_id"`
	ProjectID      string `json:"project_id"`
	ProjectRoundID string `json:"project_round_id"`
}

type UpdateGroupRequest struct {
	GroupID   string              `json:"group_id"`
	ProjectID *string             `json:"project_id,omitempty"`
	LeaderID  *string             `json:"leader_id,omitempty"`
	Setting   *model.GroupSetting `json:"setting,omitempty"`
}

type CreateJoinGroupRequest struct {
	GroupID string `json:"group_id"`
	Message string `json:"message"`
}

type UpdateJoinGroupRequest struct {
	GroupID   string                       `json:"group_id"`
	RequestID string                       `json:"request_id"`
	Status    model.JoinGroupRequestStatus `json:"status"`
}

type CreateGroupInvitationRequest struct {
	GroupID     string `json:"group_id"`
	RecipientID string `json:"recipient_id"`
}

type UpdateGroupInvitationRequest struct {
	GroupID      string                       `json:"group_id"`
	InvitationID string                       `json:"invitation_id"`
	Status       model.JoinGroupRequestStatus `json:"status"`
}

type CreateTaskRequest struct {
	GroupID     string    `json:"group_id"`
	Title       string    `json:"title"`
	Details     *string   `json:"details,omitempty"`
	AssignToIDs []string  `json:"assign_to_ids"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

type UpdateTaskRequest struct {
	GroupID     string     `json:"group_id"`
	TaskID      string     `json:"task_id"`
	Title       *string    `json:"title,omitempty"`
	Details     *string    `json:"details,omitempty"`
	AssignToIDs []string   `json:"assign_to_ids,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      *string    `json:"status,omitempty"`
}

type CreateReportRequest struct {
	ClassroomID    string             `json:"classroom_id"`
	GroupID        string             `json:"group_id"`
	ProjectRoundID string             `json:"project_round_id"`
	ReportPeriodID string             `json:"report_period_id"`
	Title          string             `json:"title"`
	Content        string             `json:"content"`
	Attachments    []AttachmentUpload `json:"attachments,omitempty"`
}

type UpdateReportRequest struct {
	GroupID          string             `json:"group_id"`
	ReportID         string             `json:"report_id"`
	Title            *string            `json:"title,omitempty"`
	Content          *string            `json:"content,omitempty"`
	AttachmentsToAdd []AttachmentUpload `json:"attachments_to_add,omitempty"`
	FilesToRemove    []string           `json:"files_to_remove,omitempty"` // file URLs to remove
}

type CreateReportFeedbackRequest struct {
	GroupID  string `json:"group_id"`
	ReportID string `json:"report_id"`
	Content  string `json:"content"`
	Grade    string `json:"grade"`
}

type UpdateReportFeedbackRequest struct {
	GroupID    string  `json:"group_id"`
	FeedbackID string  `json:"feedback_id"`
	Content    *string `json:"content,omitempty"`
	Grade      *string `json:"grade,omitempty"`
}

type GetGroupsFilterQuery struct {
	ClassroomID string  `form:"classroom_id"`
	ProjectID   *string `form:"project_id,omitempty"`
	MemberID    *string `form:"member_id,omitempty"`
}

type GroupResponse struct {
	ID             string             `json:"id"`
	ClassroomID    string             `json:"classroom_id"`
	ProjectID      string             `json:"project_id"`
	GroupChannelID string             `json:"group_channel_id"`
	LeaderID       string             `json:"leader_id"`
	Members        []UserInfoResponse `json:"members"`
	Tasks          []TaskResponse     `json:"tasks"`
	TaskStatuses   []string           `json:"task_statuses"`
	Reports        []ReportResponse   `json:"reports"`
	Setting        model.GroupSetting `json:"setting"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	AssignToIDs []string  `json:"assign_to_ids"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

type ReportResponse struct {
	ID             string                 `json:"id"`
	ReportPeriodID string                 `json:"report_period_id"`
	Title          string                 `json:"title"`
	Content        string                 `json:"content"`
	Attachments    []model.Attachment     `json:"attachments,omitempty"`
	Feedback       ReportFeedbackResponse `json:"feedback"`
}

type ReportFeedbackResponse struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	Grade       string    `json:"grade"`
	LecturerID  string    `json:"lecturer_id"`
	CommentedAt time.Time `json:"commented_at"`
}

func FromGroup(group *model.Group) *GroupResponse {
	if group == nil {
		return nil
	}

	// Convert Tasks
	tasks := make([]TaskResponse, len(group.Tasks))
	for i, t := range group.Tasks {
		assignToIDs := make([]string, len(t.AssignToIDs))
		for j, oid := range t.AssignToIDs {
			assignToIDs[j] = oid.Hex()
		}

		tasks[i] = TaskResponse{
			ID:          t.ID.Hex(),
			Title:       t.Title,
			AssignToIDs: assignToIDs,
			DueDate:     t.DueDate,
			Status:      t.Status,
		}
	}

	// Convert Reports
	reports := make([]ReportResponse, len(group.Reports))
	for i, r := range group.Reports {
		reports[i] = ReportResponse{
			ID:             r.ID.Hex(),
			ReportPeriodID: r.ReportPeriodID.Hex(),
			Title:          r.Title,
			Content:        r.Content,
			Attachments:    r.Attachments,
			Feedback: ReportFeedbackResponse{
				ID:          r.Feedback.ID.Hex(),
				Content:     r.Feedback.Content,
				Grade:       r.Feedback.Grade,
				LecturerID:  r.Feedback.LecturerID.Hex(),
				CommentedAt: r.Feedback.CommentedAt,
			},
		}
	}

	// Convert Members
	members := make([]UserInfoResponse, len(group.Members))
	for i, m := range group.Members {
		members[i] = UserInfoResponse{
			UserID:      m.ID.Hex(),
			FullName:    m.FullName,
			Email:       m.Email,
			Avatar:      m.Avatar,
			StudentCode: m.StudentCode,
		}
	}

	return &GroupResponse{
		ID:             group.ID.Hex(),
		ClassroomID:    group.ClassroomID.Hex(),
		ProjectID:      group.ProjectID.Hex(),
		GroupChannelID: group.GroupChannelID.Hex(),
		LeaderID:       group.LeaderID.Hex(),
		Members:        members,
		Tasks:          tasks,
		TaskStatuses:   group.TaskStatuses,
		Reports:        reports,
		Setting:        group.Setting,
	}
}

func FromGroups(groups []model.Group) []GroupResponse {
	if len(groups) == 0 {
		return []GroupResponse{}
	}

	responses := make([]GroupResponse, len(groups))
	for i, g := range groups {
		gr := FromGroup(&g)
		if gr != nil {
			responses[i] = *gr
		}
	}

	return responses
}
