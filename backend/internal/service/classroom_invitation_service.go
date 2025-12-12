package service

import (
	"errors"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/email"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassroomInvitationService interface {
	// Invite members to classroom
	InviteMembers(classroomID, inviterID string, emails []string) dto.InviteResult

	// Get pending invitations
	GetMyPendingInvitations(userID string) ([]dto.InvitationResponse, error)
	GetClassroomInvitations(classroomID string, status *model.InvitationStatus, page, pageSize int) ([]dto.InvitationResponse, int64, error)

	// Accept/Reject invitation
	AcceptInvitation(invitationID, userID string) (dto.AcceptInvitationResponse, error)
	RejectInvitation(invitationID, userID string) error

	// Cancel invitation (by classroom owner)
	CancelInvitation(invitationID, requesterID string) error
}

type classroomInvitationService struct {
	invitationRepo repo.ClassroomInvitationRepo
	classroomRepo  repo.ClassroomRepo
	userRepo       repo.UserRepo
	emailSender    email.Sender
}

func NewClassroomInvitationService(
	invitationRepo repo.ClassroomInvitationRepo,
	classroomRepo repo.ClassroomRepo,
	userRepo repo.UserRepo,
	emailSender email.Sender,
) ClassroomInvitationService {
	return &classroomInvitationService{
		invitationRepo: invitationRepo,
		classroomRepo:  classroomRepo,
		userRepo:       userRepo,
		emailSender:    emailSender,
	}
}

// InviteMembers invites multiple emails to a classroom
func (s *classroomInvitationService) InviteMembers(classroomID, inviterID string, emails []string) dto.InviteResult {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	result := dto.InviteResult{
		Invited:        make([]string, 0),
		AlreadyMembers: make([]string, 0),
		Failed:         make([]dto.InviteError, 0),
		TotalProcessed: len(emails),
	}

	// Get classroom info
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		// If classroom not found, all emails fail
		for _, email := range emails {
			result.Failed = append(result.Failed, dto.InviteError{
				Email: email,
				Error: "Classroom not found",
			})
		}
		return result
	}

	// Get inviter info
	inviterObjectID, err := primitive.ObjectIDFromHex(inviterID)
	if err != nil {
		for _, email := range emails {
			result.Failed = append(result.Failed, dto.InviteError{
				Email: email,
				Error: "Invalid inviter ID",
			})
		}
		return result
	}

	// Process each email
	for _, email := range emails {
		if email == "" {
			result.Failed = append(result.Failed, dto.InviteError{
				Email: email,
				Error: "Empty email",
			})
			continue
		}

		// Check if user already exists
		user, err := s.userRepo.GetByEmail(ctx, email)

		var invitedToID *primitive.ObjectID
		if err == nil {
			// User exists
			invitedToID = &user.ID

			// Check if already in classroom
			isInClassroom, err := s.classroomRepo.IsStudentInClassroom(ctx, classroomID, user.ID.Hex())
			if err == nil && isInClassroom {
				result.AlreadyMembers = append(result.AlreadyMembers, email)
				continue
			}
		}

		// Check if invitation already exists (pending)
		existingInvitation, err := s.invitationRepo.GetByEmailAndClassroom(ctx, email, classroomID)
		if err == nil && existingInvitation.Status == model.InvitationStatusPending {
			// Already invited, resend email
			go s.sendInvitationEmail(email, classroom.ID.Hex(), inviterID)
			result.Invited = append(result.Invited, email)
			continue
		}

		// Create new invitation
		invitation := &model.ClassroomInvitation{
			Email:       email,
			ClassroomID: classroom.ID,
			InvitedBy:   inviterObjectID,
			InvitedTo:   invitedToID,
			Status:      model.InvitationStatusPending,
			ExpiresAt:   time.Now().Add(30 * 24 * time.Hour), // 30 days
			CreatedAt:   time.Now(),
		}

		_, err = s.invitationRepo.Create(ctx, invitation)
		if err != nil {
			result.Failed = append(result.Failed, dto.InviteError{
				Email: email,
				Error: "Failed to create invitation",
			})
			continue
		}

		// Send invitation email
		go s.sendInvitationEmail(email, classroom.ID.Hex(), inviterID)

		result.Invited = append(result.Invited, email)
	}

	return result
}

// GetMyPendingInvitations gets all pending invitations for current user
func (s *classroomInvitationService) GetMyPendingInvitations(userID string) ([]dto.InvitationResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get user email
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	// Get pending invitations by email
	invitations, err := s.invitationRepo.GetPendingByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs with classroom names
	responses := make([]dto.InvitationResponse, 0, len(invitations))
	for _, inv := range invitations {
		response := dto.FromInvitation(&inv)

		// Get classroom name
		classroom, err := s.classroomRepo.GetByID(ctx, inv.ClassroomID.Hex())
		if err == nil {
			response.ClassroomName = classroom.UniversityName // TODO: Fix this - should be classroom name
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// GetClassroomInvitations gets invitations for a classroom (for owner/admin)
func (s *classroomInvitationService) GetClassroomInvitations(classroomID string, status *model.InvitationStatus, page, pageSize int) ([]dto.InvitationResponse, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	invitations, total, err := s.invitationRepo.GetByClassroom(ctx, classroomID, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return dto.FromInvitations(invitations), total, nil
}

// AcceptInvitation accepts an invitation and adds user to classroom
func (s *classroomInvitationService) AcceptInvitation(invitationID, userID string) (dto.AcceptInvitationResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get invitation
	invitation, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return dto.AcceptInvitationResponse{}, err
	}

	// Validate invitation status
	if invitation.Status != model.InvitationStatusPending {
		return dto.AcceptInvitationResponse{}, apperror.ErrInvitationAlreadyProcessed
	}

	// Check if expired
	if invitation.ExpiresAt.Before(time.Now()) {
		return dto.AcceptInvitationResponse{}, apperror.ErrInvitationExpired
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.AcceptInvitationResponse{}, apperror.ErrUserNotFound
		}
		return dto.AcceptInvitationResponse{}, err
	}

	// Verify email matches
	if invitation.Email != user.Email {
		return dto.AcceptInvitationResponse{}, apperror.ErrEmailMismatch
	}

	// Check if already in classroom
	isInClassroom, err := s.classroomRepo.IsStudentInClassroom(ctx, invitation.ClassroomID.Hex(), userID)
	if err == nil && isInClassroom {
		return dto.AcceptInvitationResponse{}, apperror.ErrAlreadyInClassroom
	}

	// Add user to classroom
	studentInfo := model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   "",
	}
	if user.Avatar != nil {
		studentInfo.Avatar = user.Avatar.URL
	}

	err = s.classroomRepo.AddStudent(ctx, invitation.ClassroomID.Hex(), studentInfo)
	if err != nil {
		return dto.AcceptInvitationResponse{}, err
	}

	// Update invitation status
	err = s.invitationRepo.UpdateStatus(ctx, invitationID, model.InvitationStatusAccepted)
	if err != nil {
		return dto.AcceptInvitationResponse{}, err
	}

	// Get classroom info for response
	classroom, err := s.classroomRepo.GetByID(ctx, invitation.ClassroomID.Hex())
	if err != nil {
		return dto.AcceptInvitationResponse{}, err
	}

	return dto.AcceptInvitationResponse{
		ClassroomID:   classroom.ID.Hex(),
		ClassroomName: classroom.UniversityName, // TODO: Fix - should be classroom name
	}, nil
}

// RejectInvitation rejects an invitation
func (s *classroomInvitationService) RejectInvitation(invitationID, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get invitation
	invitation, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}

	// Validate invitation status
	if invitation.Status != model.InvitationStatusPending {
		return apperror.ErrInvitationAlreadyProcessed
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	// Verify email matches
	if invitation.Email != user.Email {
		return apperror.ErrEmailMismatch
	}

	// Update invitation status to rejected
	err = s.invitationRepo.UpdateStatus(ctx, invitationID, model.InvitationStatusRejected)
	if err != nil {
		return err
	}

	return nil
}

// CancelInvitation cancels an invitation (by classroom owner)
func (s *classroomInvitationService) CancelInvitation(invitationID, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get invitation
	invitation, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}

	// Get classroom to check if requester is owner/lecturer
	classroom, err := s.classroomRepo.GetByID(ctx, invitation.ClassroomID.Hex())
	if err != nil {
		return err
	}

	// Check if requester is the lecturer
	if classroom.Lecturer.ID.Hex() != requesterID {
		return apperror.ErrForbidden
	}

	// Validate invitation status
	if invitation.Status != model.InvitationStatusPending {
		return apperror.ErrInvitationAlreadyProcessed
	}

	// Update invitation status to canceled
	err = s.invitationRepo.UpdateStatus(ctx, invitationID, model.InvitationStatusCanceled)
	if err != nil {
		return err
	}

	return nil
}

// sendInvitationEmail sends invitation email to user
func (s *classroomInvitationService) sendInvitationEmail(email, classroomID, inviterID string) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get classroom info
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return
	}

	// Get inviter info
	inviter, err := s.userRepo.GetByID(ctx, inviterID)
	if err != nil {
		return
	}

	// Send email
	classroomName := classroom.UniversityName // TODO: Fix - should be classroom name
	inviterName := inviter.Username

	err = s.emailSender.SendClassroomInvitation(email, classroomName, inviterName)
	if err != nil {
		// Log error but don't fail the invitation
		// TODO: Add proper logging
	}
}
