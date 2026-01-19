package dto

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

type CreateProjectRoundRequest struct {
	ClassroomID string `json:"classroom_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	StartDate   string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate     string `json:"end_date" binding:"required,datetime=2006-01-02"`
	Description string `json:"description"`
}

type CreateProjectRoundsRequest struct {
	ClassroomID string                      `json:"classroom_id" binding:"required"`
	Rounds      []CreateProjectRoundRequest `json:"rounds" binding:"required,min=1"`
}
type UpdateProjectRoundRequest struct {
	ClassroomID    string `json:"classroom_id" binding:"required"`
	ProjectRoundID string `json:"project_round_id" binding:"required"`
	Name           string `json:"name,omitempty"`
	StartDate      string `json:"start_date,omitempty" binding:"datetime=2006-01-02"`
	EndDate        string `json:"end_date,omitempty" binding:"datetime=2006-01-02"`
	Description    string `json:"description,omitempty"`
}

type CreateProjectRequest struct {
	ClassroomID    string `json:"classroom_id" binding:"required"`
	ProjectRoundID string `json:"project_round_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Amount         int    `json:"amount" binding:"required,min=1"`
	Description    string `json:"description"`
	MinMember      int    `json:"min_member" binding:"required,min=1"`
	MaxMember      int    `json:"max_member" binding:"required,min=1"`
}
type CreateProjectsRequest struct {
	ClassroomID    string                 `json:"classroom_id" binding:"required"`
	ProjectRoundID string                 `json:"project_round_id" binding:"required"`
	Projects       []CreateProjectRequest `json:"projects" binding:"required,min=1"`
}

type UpdateProjectRequest struct {
	ClassroomID string `json:"classroom_id" binding:"required"`
	ProjectID   string `json:"project_id" binding:"required"`
	Title       string `json:"title,omitempty"`
	Amount      int    `json:"amount,omitempty" binding:"min=1"`
	Description string `json:"description,omitempty"`
	MinMember   int    `json:"min_member,omitempty" binding:"min=1"`
	MaxMember   int    `json:"max_member,omitempty" binding:"min=1"`
}

type CreateReportPeriodRequest struct {
	ClassroomID    string   `json:"classroom_id" binding:"required"`
	ProjectRoundID string   `json:"project_round_id" binding:"required"`
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description"`
	FileType       []string `json:"file_type" binding:"required,dive,required"`
	StartDate      string   `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate        string   `json:"end_date" binding:"required,datetime=2006-01-02"`
}

type CreateReportPeriodsRequest struct {
	ClassroomID    string                      `json:"classroom_id" binding:"required"`
	ProjectRoundID string                      `json:"project_round_id" binding:"required"`
	ReportPeriods  []CreateReportPeriodRequest `json:"report_periods" binding:"required,min=1"`
}

type UpdateReportPeriodRequest struct {
	ClassroomID    string   `json:"classroom_id" binding:"required"`
	ProjectRoundID string   `json:"project_round_id" binding:"required"`
	ReportPeriodID string   `json:"report_period_id" binding:"required"`
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	FileType       []string `json:"file_type,omitempty" binding:"dive,required"`
	StartDate      string   `json:"start_date,omitempty" binding:"datetime=2006-01-02"`
	EndDate        string   `json:"end_date,omitempty" binding:"datetime=2006-01-02"`
}

// Response DTOs

type ProjectRoundResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	StartDate     time.Time              `json:"start_date"`
	EndDate       time.Time              `json:"end_date"`
	Description   string                 `json:"description"`
	ReportPeriods []ReportPeriodResponse `json:"report_periods"`
	CreatedAt     time.Time              `json:"created_at"`
	IsDeleted     bool                   `json:"is_deleted"`
}

type ProjectResponse struct {
	ID             string `json:"id"`
	ClassroomID    string `json:"classroom_id"`
	ProjectRoundID string `json:"project_round_id"`
	Title          string `json:"title"`
	Amount         int    `json:"amount"`
	Description    string `json:"description"`
	MinMember      int    `json:"min_member"`
	MaxMember      int    `json:"max_member"`
	Status         string `json:"status"`
}

type ReportPeriodResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileType    []string  `json:"file_type"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type ReportPeriodWithProjectRound struct {
	ClassroomID    string             `json:"classroom_id"`
	ProjectRoundID string             `json:"project_round_id"`
	ReportPeriods  model.ReportPeriod `json:"report_periods"`
}

type ProjectRoundWithClassroom struct {
	ClassroomID   string             `json:"classroom_id"`
	ClassroomName string             `json:"classroom_name"`
	Round         model.ProjectRound `json:"round"`
}

// Converter functions

func FromProjectRound(round *model.ProjectRound) ProjectRoundResponse {
	if round == nil {
		return ProjectRoundResponse{}
	}

	reportPeriods := make([]ReportPeriodResponse, 0, len(round.ReportPeriods))
	for _, rp := range round.ReportPeriods {
		reportPeriods = append(reportPeriods, FromReportPeriod(&rp))
	}

	return ProjectRoundResponse{
		ID:            round.ID.Hex(),
		Name:          round.Name,
		StartDate:     round.StartDate,
		EndDate:       round.EndDate,
		Description:   round.Description,
		ReportPeriods: reportPeriods,
		CreatedAt:     round.CreatedAt,
		IsDeleted:     round.IsDeleted,
	}
}

func FromProjectRounds(rounds []model.ProjectRound) []ProjectRoundResponse {
	result := make([]ProjectRoundResponse, 0, len(rounds))
	for _, r := range rounds {
		result = append(result, FromProjectRound(&r))
	}
	return result
}

func FromProject(project *model.Project) ProjectResponse {
	if project == nil {
		return ProjectResponse{}
	}

	return ProjectResponse{
		ID:             project.ID.Hex(),
		ClassroomID:    project.ClassroomID.Hex(),
		ProjectRoundID: project.ProjectRoundID.Hex(),
		Title:          project.Title,
		Amount:         project.Amount,
		Description:    project.Description,
		MinMember:      project.MinMember,
		MaxMember:      project.MaxMember,
		Status:         string(project.Status),
	}
}

func FromProjects(projects []model.Project) []ProjectResponse {
	result := make([]ProjectResponse, 0, len(projects))
	for _, p := range projects {
		result = append(result, FromProject(&p))
	}
	return result
}

func FromReportPeriod(period *model.ReportPeriod) ReportPeriodResponse {
	if period == nil {
		return ReportPeriodResponse{}
	}

	return ReportPeriodResponse{
		ID:          period.ID.Hex(),
		Title:       period.Title,
		Description: period.Description,
		FileType:    period.FileType,
		StartDate:   period.StartDate,
		EndDate:     period.EndDate,
	}
}

func FromReportPeriods(periods []model.ReportPeriod) []ReportPeriodResponse {
	result := make([]ReportPeriodResponse, 0, len(periods))
	for _, p := range periods {
		result = append(result, FromReportPeriod(&p))
	}
	return result
}
