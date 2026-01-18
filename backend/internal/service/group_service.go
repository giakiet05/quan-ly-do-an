package service

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupService interface {
	CreateGroup(req *dto.CreateGroupRequest, requesterID string) (*model.Group, error)
	GetGroupByID(groupID string, requesterID string) (*model.Group, error)
	GetGroupsFilter(query *dto.GetGroupsFilterQuery, requesterID string) ([]model.Group, error)
	UpdateGroup(req *dto.UpdateGroupRequest, requesterID string) (*model.Group, error)
	DeleteGroup(groupID string, requesterID string) error

	CreateTask(req *dto.CreateTaskRequest, requesterID string) (*model.Task, error)
	UpdateTask(req *dto.UpdateTaskRequest, requesterID string) (*model.Task, error)
	DeleteTask(groupID string, taskID string, requesterID string) error

	CreateReport(req *dto.CreateReportRequest, requesterID string) (*model.Report, error)
	UpdateReport(req *dto.UpdateReportRequest, requesterID string) (*model.Report, error)
	DeleteReport(groupID string, reportID string, requesterID string) error
	CreateReportFeedback(req *dto.CreateReportFeedbackRequest, requesterID string) (*model.ReportFeedback, error)
	UpdateReportFeedback(req *dto.UpdateReportFeedbackRequest, groupID string, reportID string) error
	DeleteReportFeedback(groupID string, reportID string, requesterID string) error
}

type groupService struct {
	groupRepo     repo.GroupRepo
	classroomRepo repo.ClassroomRepo
	channelRepo   repo.ChannelRepo
	userRepo      repo.UserRepo
}

func NewGroupService(
	groupRepo repo.GroupRepo,
	classroomRepo repo.ClassroomRepo,
	channelRepo repo.ChannelRepo,
	userRepo repo.UserRepo,
) GroupService {
	return &groupService{groupRepo: groupRepo, classroomRepo: classroomRepo, channelRepo: channelRepo, userRepo: userRepo}
}

func (g *groupService) CreateGroup(req *dto.CreateGroupRequest, requesterID string) (*model.Group, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	requesterObjectID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return nil, err
	}

	leader, err := g.userRepo.GetByID(ctx, requesterID)
	if err != nil {
		return nil, err
	}

	classroomObjectID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return nil, err
	}

	projectObjectID, err := primitive.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Initialize members slice with leader
	members := []model.UserInfo{
		{
			ID:       requesterObjectID,
			FullName: leader.FullName,
			Avatar:   leader.Avatar,
		},
	}

	// Add additional members from MemberIDs
	for _, memberIDStr := range req.MemberIDs {
		if requesterID == memberIDStr {
			continue
		}

		memberObjectID, err := primitive.ObjectIDFromHex(memberIDStr)
		if err != nil {
			return nil, apperror.ErrInvalidID
		}

		user, err := g.userRepo.GetByID(ctx, memberIDStr)
		if err != nil {
			return nil, apperror.ErrUserNotFound
		}

		members = append(members, model.UserInfo{
			ID:          memberObjectID,
			FullName:    user.FullName,
			StudentCode: user.StudentCode,
			Avatar:      user.Avatar,
		})
	}

	classroom, err := g.classroomRepo.GetByID(ctx, req.ClassroomID)
	if err != nil {
		return nil, err
	}

	// Tìm ProjectGroup theo ID
	var groupFound *model.ProjectRound
	for i, group := range classroom.ProjectRounds {
		if group.ID.Hex() == req.ProjectRoundID && !group.IsDeleted {
			groupFound = &classroom.ProjectRounds[i]
			break
		}
	}

	if groupFound == nil {
		return nil, apperror.ErrProjectGroupNotFound
	}

	// Tìm Project trong ProjectGroup
	var projectFound *model.Project
	for i, project := range groupFound.Projects {
		if project.ID.Hex() == req.ProjectID {
			projectFound = &groupFound.Projects[i]
			break
		}
	}

	if projectFound == nil {
		return nil, apperror.ErrProjectNotFound // không tìm thấy project
	}

	if len(req.MemberIDs) < projectFound.MinMember || len(req.MemberIDs) > projectFound.MaxMember {
		return nil, apperror.ErrInvalidMemberNumber
	}

	// Validate member count
	memberCount := len(members)
	if memberCount > projectFound.MaxMember || memberCount < projectFound.MinMember {
		return nil, apperror.ErrBadRequest
	}

	// Create group object
	group := &model.Group{
		ClassroomID: classroomObjectID,
		ProjectID:   projectObjectID,
		LeaderID:    requesterObjectID,
		Members:     members,
		Tasks:       []model.Task{},
		Reports:     []model.Report{},
		Setting:     model.GroupSetting{},
	}

	group, err = g.groupRepo.Create(ctx, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (g *groupService) GetGroupByID(groupID string, requesterID string) (*model.Group, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := g.groupRepo.IsMember(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	return g.groupRepo.GetByID(ctx, groupID)
}

func (g *groupService) GetGroupsFilter(query *dto.GetGroupsFilterQuery, requesterID string) ([]model.Group, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// If member_id is specified, check if requester is that member or is lecturer
	if query.MemberID != nil && *query.MemberID != "" {
		// Allow if requester is the member themselves
		if *query.MemberID != requesterID {
			// Or if requester is lecturer of the classroom
			ok, err := g.classroomRepo.IsLecturer(ctx, query.ClassroomID, requesterID)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, apperror.ErrForbidden
			}
		}
	} else {
		// No member_id filter, only lecturer can view all groups
		ok, err := g.classroomRepo.IsLecturer(ctx, query.ClassroomID, requesterID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, apperror.ErrForbidden
		}
	}

	groups, err := g.groupRepo.GetFilter(ctx, query.ClassroomID, query.ProjectID, query.MemberID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (g *groupService) UpdateGroup(req *dto.UpdateGroupRequest, requesterID string) (*model.Group, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is the group leader
	ok, err := g.groupRepo.IsLeader(ctx, req.GroupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	// Get existing group
	group, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// Update project if provided
	if req.ProjectID != nil {
		projectOID, err := primitive.ObjectIDFromHex(*req.ProjectID)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}
		group.ProjectID = projectOID
	}

	// Update leader if provided
	if req.LeaderID != nil {
		leaderOID, err := primitive.ObjectIDFromHex(*req.LeaderID)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}

		// Verify new leader is a member of the group
		isMember := false
		for _, member := range group.Members {
			if member.ID == leaderOID {
				isMember = true
				break
			}
		}
		if !isMember {
			return nil, apperror.ErrBadRequest
		}

		group.LeaderID = leaderOID
	}

	// Update setting if provided
	if req.Setting != nil {
		group.Setting = *req.Setting
	}

	// Replace the entire group document
	err = g.groupRepo.Replace(ctx, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (g *groupService) DeleteGroup(groupID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := g.groupRepo.IsLeader(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return g.groupRepo.Delete(ctx, groupID)
}

func (g *groupService) CreateTask(req *dto.CreateTaskRequest, requesterID string) (*model.Task, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is the group leader
	ok, err := g.groupRepo.IsLeader(ctx, req.GroupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	// Get existing group to verify members
	group, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// Convert assign to IDs
	assignToOIDs := make([]primitive.ObjectID, 0, len(req.AssignToIDs))
	for _, assignID := range req.AssignToIDs {
		oid, err := primitive.ObjectIDFromHex(assignID)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}

		// Verify assigned user is a member
		isMember := false
		for _, member := range group.Members {
			if member.ID == oid {
				isMember = true
				break
			}
		}
		if !isMember {
			return nil, apperror.ErrBadRequest
		}

		assignToOIDs = append(assignToOIDs, oid)
	}

	// Create new task
	task := model.Task{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Details:     req.Details,
		AssignToIDs: assignToOIDs,
		DueDate:     req.DueDate,
		Status:      "pending",
	}

	// Update group with new task
	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id": groupOID,
	}

	update := repo.UpdateDocument{
		"$push": bson.M{
			"tasks": task,
		},
	}

	err = g.groupRepo.Update(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (g *groupService) UpdateTask(req *dto.UpdateTaskRequest, requesterID string) (*model.Task, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is the group leader
	ok, err := g.groupRepo.IsLeader(ctx, req.GroupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	// Get existing group to verify task exists and members
	group, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// Find the task to verify it exists
	taskOID, err := primitive.ObjectIDFromHex(req.TaskID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	taskIndex := -1
	for i, task := range group.Tasks {
		if task.ID == taskOID {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return nil, apperror.ErrNotFound
	}

	// Build update document
	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id":       groupOID,
		"tasks._id": taskOID,
	}

	setFields := bson.M{}

	// Update title if provided
	if req.Title != nil {
		setFields["tasks.$.title"] = *req.Title
	}

	// Update details if provided
	if req.Details != nil {
		setFields["tasks.$.details"] = req.Details
	}

	// Update assign to IDs if provided
	if req.AssignToIDs != nil {
		assignToOIDs := make([]primitive.ObjectID, 0, len(req.AssignToIDs))
		for _, assignID := range req.AssignToIDs {
			oid, err := primitive.ObjectIDFromHex(assignID)
			if err != nil {
				return nil, apperror.ErrBadRequest
			}

			// Verify assigned user is a member
			isMember := false
			for _, member := range group.Members {
				if member.ID == oid {
					isMember = true
					break
				}
			}
			if !isMember {
				return nil, apperror.ErrBadRequest
			}

			assignToOIDs = append(assignToOIDs, oid)
		}
		setFields["tasks.$.assign_to_ids"] = assignToOIDs
	}

	// Update due date if provided
	if req.DueDate != nil {
		setFields["tasks.$.due_date"] = *req.DueDate
	}

	// Update status if provided
	if req.Status != nil {
		setFields["tasks.$.status"] = *req.Status
	}

	if len(setFields) > 0 {
		update := repo.UpdateDocument{
			"$set": setFields,
		}

		err = g.groupRepo.Update(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	// Get updated task
	updatedGroup, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	for _, task := range updatedGroup.Tasks {
		if task.ID == taskOID {
			return &task, nil
		}
	}

	return nil, apperror.ErrNotFound
}

func (g *groupService) DeleteTask(groupID string, taskID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is the group leader
	ok, err := g.groupRepo.IsLeader(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	// Verify task exists
	group, err := g.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}

	taskOID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	found := false
	for _, task := range group.Tasks {
		if task.ID == taskOID {
			found = true
			break
		}
	}

	if !found {
		return apperror.ErrNotFound
	}

	// Delete task
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id": groupOID,
	}

	update := repo.UpdateDocument{
		"$pull": bson.M{
			"tasks": bson.M{
				"_id": taskOID,
			},
		},
	}

	return g.groupRepo.Update(ctx, filter, update)
}

func (g *groupService) CreateReport(req *dto.CreateReportRequest, requesterID string) (*model.Report, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	reportPeriodOID, err := primitive.ObjectIDFromHex(req.ReportPeriodID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	// Check if requester is a member of the group
	ok, err := g.groupRepo.IsMember(ctx, req.GroupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	// TODO Check if report for the period already exists and if report period is valid

	now := time.Now()
	report := model.Report{
		ReportPeriodID: reportPeriodOID,
		Title:          req.Title,
		Content:        req.Content,
		Files:          req.Files,
		Feedback:       model.ReportFeedback{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id": groupOID,
	}

	update := repo.UpdateDocument{
		"$push": bson.M{
			"reports": report,
		},
	}

	err = g.groupRepo.Update(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (g *groupService) UpdateReport(req *dto.UpdateReportRequest, requesterID string) (*model.Report, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is a member of the group
	ok, err := g.groupRepo.IsMember(ctx, req.GroupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	// Get existing group to verify report exists
	group, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// Find the report to verify it exists
	reportOID, err := primitive.ObjectIDFromHex(req.ReportID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	reportIndex := -1
	for i, report := range group.Reports {
		if report.ID == reportOID {
			reportIndex = i
			break
		}
	}

	if reportIndex == -1 {
		return nil, apperror.ErrNotFound
	}

	// Build update document
	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id":         groupOID,
		"reports._id": reportOID,
	}

	setFields := bson.M{
		"reports.$.updated_at": time.Now(),
	}

	// Update title if provided
	if req.Title != nil {
		setFields["reports.$.title"] = *req.Title
	}

	// Update content if provided
	if req.Content != nil {
		setFields["reports.$.content"] = *req.Content
	}

	// Update files if provided
	if req.Files != nil {
		setFields["reports.$.files"] = req.Files
	}

	update := repo.UpdateDocument{
		"$set": setFields,
	}

	err = g.groupRepo.Update(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Get updated report
	updatedGroup, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	for _, report := range updatedGroup.Reports {
		if report.ID == reportOID {
			return &report, nil
		}
	}

	return nil, apperror.ErrNotFound
}

func (g *groupService) DeleteReport(groupID string, reportID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is the group leader
	ok, err := g.groupRepo.IsLeader(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	// Verify report exists
	group, err := g.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}

	reportOID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	found := false
	for _, report := range group.Reports {
		if report.ID == reportOID {
			found = true
			break
		}
	}

	if !found {
		return apperror.ErrNotFound
	}

	// Delete report
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id": groupOID,
	}

	update := repo.UpdateDocument{
		"$pull": bson.M{
			"reports": bson.M{
				"_id": reportOID,
			},
		},
	}

	return g.groupRepo.Update(ctx, filter, update)
}

func (g *groupService) CreateReportFeedback(req *dto.CreateReportFeedbackRequest, requesterID string) (*model.ReportFeedback, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	group, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	ok, err := g.classroomRepo.IsLecturer(ctx, group.ClassroomID.Hex(), requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	requesterOID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	// Verify report exists
	reportOID, err := primitive.ObjectIDFromHex(req.ReportID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	found := false
	for _, report := range group.Reports {
		if report.ID == reportOID {
			found = true
			break
		}
	}

	if !found {
		return nil, apperror.ErrNotFound
	}

	// Create feedback
	now := time.Now()
	feedback := model.ReportFeedback{
		Content:     req.Content,
		Grade:       req.Grade,
		LecturerID:  requesterOID,
		CommentedAt: now,
		UpdatedAt:   now,
	}

	// Update report with feedback
	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id":         groupOID,
		"reports._id": reportOID,
	}

	update := repo.UpdateDocument{
		"$set": bson.M{
			"reports.$.feedback": feedback,
		},
	}

	err = g.groupRepo.Update(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &feedback, nil
}

func (g *groupService) UpdateReportFeedback(req *dto.UpdateReportFeedbackRequest, groupID string, reportID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get the group to verify report exists
	group, err := g.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}

	// Find the report to verify it exists and has feedback
	reportOID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	found := false
	hasFeedback := false
	for _, report := range group.Reports {
		if report.ID == reportOID {
			found = true
			if !report.Feedback.LecturerID.IsZero() {
				hasFeedback = true
			}
			break
		}
	}

	if !found {
		return apperror.ErrNotFound
	}

	if !hasFeedback {
		return apperror.ErrNotFound
	}

	// Build update document
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id":         groupOID,
		"reports._id": reportOID,
	}

	setFields := bson.M{
		"reports.$.feedback.updated_at": time.Now(),
	}

	// Update content if provided
	if req.Content != nil {
		setFields["reports.$.feedback.content"] = *req.Content
	}

	// Update grade if provided
	if req.Grade != nil {
		setFields["reports.$.feedback.grade"] = *req.Grade
	}

	update := repo.UpdateDocument{
		"$set": setFields,
	}

	return g.groupRepo.Update(ctx, filter, update)
}

func (g *groupService) DeleteReportFeedback(groupID string, reportID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get the group
	group, err := g.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}

	ok, err := g.classroomRepo.IsLecturer(ctx, group.ClassroomID.Hex(), requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	// Verify report exists and has feedback
	reportOID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	found := false
	hasFeedback := false
	for _, report := range group.Reports {
		if report.ID == reportOID {
			found = true
			if !report.Feedback.LecturerID.IsZero() {
				hasFeedback = true
			}
			break
		}
	}

	if !found {
		return apperror.ErrNotFound
	}

	if !hasFeedback {
		return apperror.ErrNotFound
	}

	// Clear feedback
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return apperror.ErrBadRequest
	}

	filter := repo.Filter{
		"_id":         groupOID,
		"reports._id": reportOID,
	}

	update := repo.UpdateDocument{
		"$set": bson.M{
			"reports.$.feedback": model.ReportFeedback{},
		},
	}

	return g.groupRepo.Update(ctx, filter, update)
}
