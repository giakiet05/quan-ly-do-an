package service

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/bus"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupService interface {
	CreateGroup(req *dto.CreateGroupRequest, requesterID string) (*model.Group, error)
	GetGroupByID(groupID string, requesterID string) (*model.Group, error)
	GetGroupsFilter(query *dto.GetGroupsFilterQuery, requesterID string) ([]model.Group, error)
	UpdateGroup(req *dto.UpdateGroupRequest, requesterID string) (*model.Group, error)
	DeleteGroup(groupID string, requesterID string) error

	// Join requests (student requests to join a group)
	CreateJoinRequest(req *dto.CreateJoinGroupRequest, requesterID string) (*model.JoinGroupRequest, error)
	AcceptJoinRequest(req *dto.UpdateJoinGroupRequest, requesterID string) error
	RejectJoinRequest(req *dto.UpdateJoinGroupRequest, requesterID string) error

	// Invitations (group leader invites a student)
	InviteToGroup(req *dto.CreateGroupInvitationRequest, requesterID string) (*model.JoinGroupInvitation, error)
	AcceptInvitation(req *dto.UpdateGroupInvitationRequest, requesterID string) error
	RejectInvitation(req *dto.UpdateGroupInvitationRequest, requesterID string) error

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
	projectRepo   repo.ProjectRepo
	eventBus      *bus.EventBus
	cron          *cron.Cron
}

func NewGroupService(
	groupRepo repo.GroupRepo,
	classroomRepo repo.ClassroomRepo,
	channelRepo repo.ChannelRepo,
	userRepo repo.UserRepo,
	projectRepo repo.ProjectRepo,
	eventBus *bus.EventBus,
	cron *cron.Cron,
) GroupService {
	return &groupService{
		groupRepo:     groupRepo,
		classroomRepo: classroomRepo,
		channelRepo:   channelRepo,
		userRepo:      userRepo,
		projectRepo:   projectRepo,
		eventBus:      eventBus,
		cron:          cron,
	}
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

	// Tìm ProjectRound theo ID
	var roundFound *model.ProjectRound
	for i, round := range classroom.ProjectRounds {
		if round.ID.Hex() == req.ProjectRoundID && !round.IsDeleted {
			roundFound = &classroom.ProjectRounds[i]
			break
		}
	}

	if roundFound == nil {
		return nil, apperror.ErrProjectGroupNotFound
	}

	// Tìm Project từ collection riêng
	projectFound, err := g.projectRepo.GetProjectByID(ctx, req.ClassroomID, req.ProjectID)
	if err != nil {
		return nil, err
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
			ok, err := g.classroomRepo.IsLecturerOrCoLecturer(ctx, query.ClassroomID, requesterID)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, apperror.ErrForbidden
			}
		}
	} else {
		// No member_id filter, only lecturer can view all groups
		ok, err := g.classroomRepo.IsLecturerOrCoLecturer(ctx, query.ClassroomID, requesterID)
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

// Join Requests - Student requests to join a group

func (g *groupService) CreateJoinRequest(req *dto.CreateJoinGroupRequest, requesterID string) (*model.JoinGroupRequest, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	requesterOID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	// Get group to check if join requests are allowed and user is not already a member
	group, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// Check if join requests are allowed
	if !group.Setting.AllowJoinRequest {
		return nil, apperror.ErrForbidden
	}

	// Check if user is already a member
	for _, member := range group.Members {
		if member.ID == requesterOID {
			return nil, apperror.ErrBadRequest
		}
	}

	// Check if group is at max capacity
	isMaxReached, err := g.groupRepo.IsMaxMemberReached(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if isMaxReached {
		return nil, apperror.ErrBadRequest
	}

	// Check if user already has a pending request
	for _, request := range group.JoinRequests {
		if request.UserID == requesterOID && request.Status == model.RequestPending {
			return nil, apperror.ErrBadRequest
		}
	}

	// Create join request
	joinRequest := model.JoinGroupRequest{
		ID:          primitive.NewObjectID(),
		UserID:      requesterOID,
		Status:      model.RequestPending,
		Message:     req.Message,
		RequestedAt: time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Add request to group
	filter := repo.Filter{"_id": groupOID}
	update := repo.UpdateDocument{
		"$push": bson.M{"join_requests": joinRequest},
	}

	err = g.groupRepo.Update(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &joinRequest, nil
}

func (g *groupService) AcceptJoinRequest(req *dto.UpdateJoinGroupRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Only group leader can accept join requests
	ok, err := g.groupRepo.IsLeader(ctx, req.GroupID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	requestOID, err := primitive.ObjectIDFromHex(req.RequestID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Get group and join request
	group, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return err
	}

	joinRequest, err := g.groupRepo.GetJoinRequestByID(ctx, groupOID, requestOID)
	if err != nil {
		return err
	}

	if joinRequest.Status != model.RequestPending {
		return apperror.ErrBadRequest
	}

	// Check if group is at max capacity
	isMaxReached, err := g.groupRepo.IsMaxMemberReached(ctx, req.GroupID)
	if err != nil {
		return err
	}
	if isMaxReached {
		return apperror.ErrBadRequest
	}

	// Get user info
	user, err := g.userRepo.GetByID(ctx, joinRequest.UserID.Hex())
	if err != nil {
		return err
	}

	userInfo := model.UserInfo{
		ID:          user.ID,
		FullName:    user.FullName,
		Avatar:      user.Avatar,
		StudentCode: user.StudentCode,
	}

	// Add user to group and update request status
	err = g.groupRepo.AddMember(ctx, req.GroupID, user)
	if err != nil {
		return err
	}

	// Update request status
	filter := repo.Filter{
		"_id":                groupOID,
		"join_requests._id": requestOID,
	}
	update := repo.UpdateDocument{
		"$set": bson.M{
			"join_requests.$.status":      model.RequestAccepted,
			"join_requests.$.updated_at": time.Now(),
		},
	}

	return g.groupRepo.Update(ctx, filter, update)
}

func (g *groupService) RejectJoinRequest(req *dto.UpdateJoinGroupRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Only group leader can reject join requests
	ok, err := g.groupRepo.IsLeader(ctx, req.GroupID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	requestOID, err := primitive.ObjectIDFromHex(req.RequestID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Update request status
	filter := repo.Filter{
		"_id":                groupOID,
		"join_requests._id": requestOID,
	}
	update := repo.UpdateDocument{
		"$set": bson.M{
			"join_requests.$.status":      model.RequestRejected,
			"join_requests.$.updated_at": time.Now(),
		},
	}

	return g.groupRepo.Update(ctx, filter, update)
}

// Invitations - Group leader invites a student

func (g *groupService) InviteToGroup(req *dto.CreateGroupInvitationRequest, requesterID string) (*model.JoinGroupInvitation, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Only group leader can invite
	ok, err := g.groupRepo.IsLeader(ctx, req.GroupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	recipientOID, err := primitive.ObjectIDFromHex(req.RecipientID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	// Get group
	group, err := g.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// Check if recipient is already a member
	for _, member := range group.Members {
		if member.ID == recipientOID {
			return nil, apperror.ErrBadRequest
		}
	}

	// Check if group is at max capacity
	isMaxReached, err := g.groupRepo.IsMaxMemberReached(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if isMaxReached {
		return nil, apperror.ErrBadRequest
	}

	// Check if there's already a pending invitation
	for _, invitation := range group.JoinInvitation {
		if invitation.RecipientID == recipientOID && invitation.Status == model.RequestPending {
			return nil, apperror.ErrBadRequest
		}
	}

	// Create invitation
	invitation := model.JoinGroupInvitation{
		ID:          primitive.NewObjectID(),
		GroupID:     groupOID,
		RecipientID: recipientOID,
		Status:      model.RequestPending,
		SentAt:      time.Now(),
	}

	// Add invitation to group
	filter := repo.Filter{"_id": groupOID}
	update := repo.UpdateDocument{
		"$push": bson.M{"join_invitations": invitation},
	}

	err = g.groupRepo.Update(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Publish event for notification
	g.eventBus.Publish(bus.GroupInvitationEvent{
		GroupID:     req.GroupID,
		InviterID:   requesterID,
		InviteeID:   req.RecipientID,
		IsAccepted:  false,
		SentAt:      invitation.SentAt,
		RespondedAt: nil,
	})

	return &invitation, nil
}

func (g *groupService) AcceptInvitation(req *dto.UpdateGroupInvitationRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	invitationOID, err := primitive.ObjectIDFromHex(req.InvitationID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Get invitation
	invitation, err := g.groupRepo.GetInvitationByID(ctx, groupOID, invitationOID)
	if err != nil {
		return err
	}

	// Only the invited user can accept
	if invitation.RecipientID.Hex() != requesterID {
		return apperror.ErrForbidden
	}

	if invitation.Status != model.RequestPending {
		return apperror.ErrBadRequest
	}

	// Check if group is at max capacity
	isMaxReached, err := g.groupRepo.IsMaxMemberReached(ctx, req.GroupID)
	if err != nil {
		return err
	}
	if isMaxReached {
		return apperror.ErrBadRequest
	}

	// Get user info
	user, err := g.userRepo.GetByID(ctx, requesterID)
	if err != nil {
		return err
	}

	// Add user to group
	err = g.groupRepo.AddMember(ctx, req.GroupID, user)
	if err != nil {
		return err
	}

	// Update invitation status
	now := time.Now()
	filter := repo.Filter{
		"_id":                  groupOID,
		"join_invitations._id": invitationOID,
	}
	update := repo.UpdateDocument{
		"$set": bson.M{
			"join_invitations.$.status":       model.RequestAccepted,
			"join_invitations.$.responded_at": now,
		},
	}

	return g.groupRepo.Update(ctx, filter, update)
}

func (g *groupService) RejectInvitation(req *dto.UpdateGroupInvitationRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	groupOID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	invitationOID, err := primitive.ObjectIDFromHex(req.InvitationID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Get invitation
	invitation, err := g.groupRepo.GetInvitationByID(ctx, groupOID, invitationOID)
	if err != nil {
		return err
	}

	// Only the invited user can reject
	if invitation.RecipientID.Hex() != requesterID {
		return apperror.ErrForbidden
	}

	if invitation.Status != model.RequestPending {
		return apperror.ErrBadRequest
	}

	// Update invitation status
	now := time.Now()
	filter := repo.Filter{
		"_id":                  groupOID,
		"join_invitations._id": invitationOID,
	}
	update := repo.UpdateDocument{
		"$set": bson.M{
			"join_invitations.$.status":       model.RequestRejected,
			"join_invitations.$.responded_at": now,
		},
	}

	return g.groupRepo.Update(ctx, filter, update)
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

	ok, err := g.groupRepo.IsMember(ctx, req.GroupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	ok, err = g.projectRepo.ReportPeriodExists(ctx, req.ClassroomID, req.ProjectRoundID, req.ReportPeriodID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrReportPeriodNotFound
	}

	exists, err := g.groupRepo.ReportExistsByPeriod(ctx, req.ClassroomID, req.ReportPeriodID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.ErrReportAlreadyExists
	}

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

	g.eventBus.Publish(bus.TopicReportSubmittedEvent{
		ClassroomID: req.ClassroomID,
		GroupID:     req.GroupID,
		ReportID:    report.ID.Hex(),
		SubmitterID: requesterID,
		SubmittedAt: now,
	})

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

	ok, err := g.classroomRepo.IsLecturerOrCoLecturer(ctx, group.ClassroomID.Hex(), requesterID)
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

	g.eventBus.Publish(bus.TopicReportGradedEvent{
		ClassroomID: group.ClassroomID.Hex(),
		GroupID:     req.GroupID,
		ReportID:    req.ReportID,
		GradedAt:    now,
	})

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

	ok, err := g.classroomRepo.IsLecturerOrCoLecturer(ctx, group.ClassroomID.Hex(), requesterID)
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
