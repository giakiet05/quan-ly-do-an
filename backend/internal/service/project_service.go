package service

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectService interface {
	CreateProjectRound(req *dto.CreateProjectRoundRequest, requesterID string) (*model.ProjectRound, error)
	CreateProjectRounds(req *dto.CreateProjectRoundsRequest, requesterID string) ([]model.ProjectRound, error)
	GetProjectRoundByID(classroomID, roundID string) (*model.ProjectRound, error)
	GetProjectRoundsByClassroomID(classroomID string) ([]model.ProjectRound, error)
	UpdateProjectRound(req *dto.UpdateProjectRoundRequest, requesterID string) (*model.ProjectRound, error)
	DeleteProjectRound(classroomID, roundID string, requesterID string) error

	CreateProject(req *dto.CreateProjectRequest, requesterID string) (*model.Project, error)
	CreateProjects(req *dto.CreateProjectsRequest, requesterID string) ([]model.Project, error)
	GetProjectByID(classroomID string, projectID string) (*model.Project, error)
	GetProjectsByRoundID(classroomID, roundID string) ([]model.Project, error)
	UpdateProject(req *dto.UpdateProjectRequest, requesterID string) (*model.Project, error)
	DeleteProject(classroomID, projectID string, requesterID string) error

	CreateReportPeriod(req *dto.CreateReportPeriodRequest, classroomID, roundID, requesterID string) (*model.ReportPeriod, error)
	CreateReportPeriods(req *dto.CreateReportPeriodsRequest, classroomID, requesterID string) ([]model.ReportPeriod, error)
	GetReportPeriodByID(classroomID, roundID, reportPeriodID string) (*model.ReportPeriod, error)
	GetReportPeriodsByRoundID(classroomID, roundID string) ([]model.ReportPeriod, error)
	UpdateReportPeriod(req *dto.UpdateReportPeriodRequest, requesterID string) (*model.ReportPeriod, error)
	DeleteReportPeriod(classroomID, roundID, reportPeriodID string, requesterID string) error
}

type projectService struct {
	projectRepo   repo.ProjectRepo
	classroomRepo repo.ClassroomRepo
}

func NewProjectService(projectRepo repo.ProjectRepo, classroomRepo repo.ClassroomRepo) ProjectService {
	return &projectService{
		projectRepo:   projectRepo,
		classroomRepo: classroomRepo,
	}
}

func (p *projectService) CreateProjectRound(req *dto.CreateProjectRoundRequest, requesterID string) (*model.ProjectRound, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, req.ClassroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	round := &model.ProjectRound{
		Name:          req.Name,
		StartDate:     startDate,
		EndDate:       endDate,
		Description:   req.Description,
		Projects:      []model.Project{},
		ReportPeriods: []model.ReportPeriod{},
	}

	err = p.projectRepo.CreateProjectRound(ctx, req.ClassroomID, round)
	if err != nil {
		return nil, err
	}

	return round, nil
}

func (p *projectService) CreateProjectRounds(req *dto.CreateProjectRoundsRequest, requesterID string) ([]model.ProjectRound, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, req.ClassroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	rounds := make([]model.ProjectRound, 0, len(req.Rounds))
	for _, r := range req.Rounds {
		startDate, err := time.Parse("2006-01-02", r.StartDate)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}

		endDate, err := time.Parse("2006-01-02", r.EndDate)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}

		rounds = append(rounds, model.ProjectRound{
			Name:          r.Name,
			StartDate:     startDate,
			EndDate:       endDate,
			Description:   r.Description,
			Projects:      []model.Project{},
			ReportPeriods: []model.ReportPeriod{},
		})
	}

	err = p.projectRepo.CreateProjectRounds(ctx, req.ClassroomID, rounds)
	if err != nil {
		return nil, err
	}
	return rounds, nil
}

func (p *projectService) GetProjectRoundByID(classroomID, roundID string) (*model.ProjectRound, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return p.projectRepo.GetProjectRoundByID(ctx, classroomID, roundID)
}

func (p *projectService) GetProjectRoundsByClassroomID(classroomID string) ([]model.ProjectRound, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return p.projectRepo.GetProjectRoundsByClassroomID(ctx, classroomID)
}

func (p *projectService) UpdateProjectRound(req *dto.UpdateProjectRoundRequest, requesterID string) (*model.ProjectRound, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, req.ClassroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	round, err := p.projectRepo.GetProjectRoundByID(ctx, req.ClassroomID, req.ProjectRoundID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		round.Name = req.Name
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}
		round.StartDate = startDate
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}
		round.EndDate = endDate
	}
	if req.Description != "" {
		round.Description = req.Description
	}

	return round, nil
}

func (p *projectService) DeleteProjectRound(classroomID, roundID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, requesterID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	return p.projectRepo.DeleteProjectRound(ctx, classroomID, roundID)
}

func (p *projectService) CreateProject(req *dto.CreateProjectRequest, requesterID string) (*model.Project, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, req.ClassroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	classroomOID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	roundOID, err := primitive.ObjectIDFromHex(req.ProjectRoundID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	project := &model.Project{
		ID:             primitive.NewObjectID(),
		ClassroomID:    classroomOID,
		ProjectRoundID: roundOID,
		Title:          req.Title,
		Amount:         req.Amount,
		Description:    req.Description,
		MinMember:      req.MinMember,
		MaxMember:      req.MaxMember,
		Status:         model.ProjectStatusPending,
	}

	err = p.projectRepo.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (p *projectService) CreateProjects(req *dto.CreateProjectsRequest, requesterID string) ([]model.Project, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, req.ClassroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	classroomOID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	roundOID, err := primitive.ObjectIDFromHex(req.ProjectRoundID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	projects := make([]model.Project, 0, len(req.Projects))
	for _, p := range req.Projects {
		projects = append(projects, model.Project{
			ID:             primitive.NewObjectID(),
			ClassroomID:    classroomOID,
			ProjectRoundID: roundOID,
			Title:          p.Title,
			Amount:         p.Amount,
			Description:    p.Description,
			MinMember:      p.MinMember,
			MaxMember:      p.MaxMember,
			Status:         model.ProjectStatusPending,
		})
	}

	err = p.projectRepo.CreateProjects(ctx, projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (p *projectService) GetProjectByID(classroomID string, projectID string) (*model.Project, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return p.projectRepo.GetProjectByID(ctx, classroomID, projectID)
}

func (p *projectService) GetProjectsByRoundID(classroomID, roundID string) ([]model.Project, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return p.projectRepo.GetProjectsByRoundID(ctx, classroomID, roundID)
}

func (p *projectService) UpdateProject(req *dto.UpdateProjectRequest, requesterID string) (*model.Project, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, req.ClassroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	project, err := p.projectRepo.GetProjectByID(ctx, req.ClassroomID, req.ProjectID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		project.Title = req.Title
	}
	if req.Amount > 0 {
		project.Amount = req.Amount
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.MinMember > 0 {
		project.MinMember = req.MinMember
	}
	if req.MaxMember > 0 {
		project.MaxMember = req.MaxMember
	}

	err = p.projectRepo.ReplaceProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (p *projectService) DeleteProject(classroomID, projectID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, requesterID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	return p.projectRepo.DeleteProject(ctx, classroomID, projectID)
}

func (p *projectService) CreateReportPeriod(req *dto.CreateReportPeriodRequest, classroomID, roundID, requesterID string) (*model.ReportPeriod, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	period := &model.ReportPeriod{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		FileType:    req.FileType,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err = p.projectRepo.CreateReportPeriod(ctx, classroomID, roundID, period)
	if err != nil {
		return nil, err
	}

	return period, nil
}

func (p *projectService) CreateReportPeriods(req *dto.CreateReportPeriodsRequest, classroomID, requesterID string) ([]model.ReportPeriod, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	periods := make([]model.ReportPeriod, 0, len(req.ReportPeriods))
	for _, p := range req.ReportPeriods {
		startDate, err := time.Parse("2006-01-02", p.StartDate)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}

		endDate, err := time.Parse("2006-01-02", p.EndDate)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}

		periods = append(periods, model.ReportPeriod{
			ID:          primitive.NewObjectID(),
			Title:       p.Title,
			Description: p.Description,
			FileType:    p.FileType,
			StartDate:   startDate,
			EndDate:     endDate,
		})
	}

	err = p.projectRepo.CreateReportPeriods(ctx, classroomID, req.ProjectRoundID, periods)
	if err != nil {
		return nil, err
	}
	return periods, nil
}

func (p *projectService) GetReportPeriodByID(classroomID, roundID, reportPeriodID string) (*model.ReportPeriod, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return p.projectRepo.GetReportPeriodByID(ctx, classroomID, roundID, reportPeriodID)
}

func (p *projectService) GetReportPeriodsByRoundID(classroomID, roundID string) ([]model.ReportPeriod, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return p.projectRepo.GetReportPeriodsByRoundID(ctx, classroomID, roundID)
}

func (p *projectService) UpdateReportPeriod(req *dto.UpdateReportPeriodRequest, requesterID string) (*model.ReportPeriod, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, req.ClassroomID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isLecturer {
		return nil, apperror.ErrForbidden
	}

	period, err := p.projectRepo.GetReportPeriodByID(ctx, req.ClassroomID, req.ProjectRoundID, req.ReportPeriodID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		period.Title = req.Title
	}
	if req.Description != "" {
		period.Description = req.Description
	}
	if len(req.FileType) > 0 {
		period.FileType = req.FileType
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}
		period.StartDate = startDate
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}
		period.EndDate = endDate
	}

	err = p.projectRepo.ReplaceReportPeriod(ctx, period)
	if err != nil {
		return nil, err
	}

	return period, nil
}

func (p *projectService) DeleteReportPeriod(classroomID, roundID, reportPeriodID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	isLecturer, err := p.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, requesterID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	return p.projectRepo.DeleteReportPeriod(ctx, classroomID, roundID, reportPeriodID)
}
