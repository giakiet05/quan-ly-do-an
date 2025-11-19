package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

type GroupResponse struct {
	ID             string                   `json:"id"`
	ClassroomID    string                   `json:"classroom_id"`
	ProjectID      string                   `json:"project_id"`
	GroupChannelID string                   `json:"group_channel_id"`
	Members        []model.UserInfoResponse `json:"members"`
	Tasks          []TaskResponse           `json:"tasks"`
	Reports        []ReportResponse         `json:"reports"`
	Setting        model.GroupSetting       `json:"setting"`
}

type TaskResponse struct {
	ID           string    `json:"id"`
	Title        string    `son:"title"`
	AssignToID   *string   `json:"assign_to_id"`
	AssignToName string    `json:"assign_to_name"`
	DueDate      time.Time `json:"due_date"`
	Status       string    `json:"status"`
}

type ReportResponse struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Files    []model.File           `json:"files"`
	Feedback ReportFeedbackResponse `json:"feedback"`
}

type ReportFeedbackResponse struct {
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
		var assignToID *string
		if t.AssignToID != nil {
			idStr := t.AssignToID.Hex()
			assignToID = &idStr
		}

		tasks[i] = TaskResponse{
			ID:           t.ID.Hex(),
			Title:        t.Title,
			AssignToID:   assignToID,
			AssignToName: t.AssignToName,
			DueDate:      t.DueDate,
			Status:       t.Status,
		}
	}

	// Convert Reports
	reports := make([]ReportResponse, len(group.Reports))
	for i, r := range group.Reports {
		reports[i] = ReportResponse{
			ID:      r.ID.Hex(),
			Title:   r.Title,
			Content: r.Content,
			Files:   r.Files,
			Feedback: ReportFeedbackResponse{
				Content:     r.Feedback.Content,
				Grade:       r.Feedback.Grade,
				LecturerID:  r.Feedback.LecturerID.Hex(),
				CommentedAt: r.Feedback.CommentedAt,
			},
		}
	}

	// Convert Members
	members := make([]model.UserInfoResponse, len(group.Members))
	for i, m := range group.Members {
		members[i] = model.UserInfoResponse{
			UserID:   m.ID.Hex(),
			Username: m.Username,
			Avatar:   m.Avatar,
		}
	}

	return &GroupResponse{
		ID:             group.ID.Hex(),
		ClassroomID:    group.ClassroomID.Hex(),
		ProjectID:      group.ProjectID.Hex(),
		GroupChannelID: group.GroupChannelID.Hex(),
		Members:        members,
		Tasks:          tasks,
		Reports:        reports,
		Setting:        group.Setting,
	}
}

func FromGroups(groups []model.Group) []GroupResponse {
	if len(groups) == 0 {
		return nil
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
