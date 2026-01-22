package service

import (
	"errors"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/platform/bus"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassroomInvitationService interface {
	InviteByEmail(classroomID, inviterID, email string) (*dto.ClassroomInvitationResponse, error)
	AcceptInvitation(invitationID, userID string) error
	RejectInvitation(invitationID, userID string) error
	CancelInvitation(invitationID, userID string) error
	GetPendingInvitationsForUser(userID string, page, pageSize int) ([]dto.ClassroomInvitationResponse, int64, error)
	GetPendingInvitationsForClassroom(classroomID, userID string, page, pageSize int) ([]dto.ClassroomInvitationResponse, int64, error)
}

type classroomInvitationService struct {
	invitationRepo repo.ClassroomInvitationRepo
	classroomRepo  repo.ClassroomRepo
	userRepo       repo.UserRepo
	channelRepo    repo.ChannelRepo
	eventBus       *bus.EventBus
}

func NewClassroomInvitationService(
	invitationRepo repo.ClassroomInvitationRepo,
	classroomRepo repo.ClassroomRepo,
	userRepo repo.UserRepo,
	channelRepo repo.ChannelRepo,
	eventBus *bus.EventBus,
) ClassroomInvitationService {
	return &classroomInvitationService{
		invitationRepo: invitationRepo,
		classroomRepo:  classroomRepo,
		userRepo:       userRepo,
		channelRepo:    channelRepo,
		eventBus:       eventBus,
	}
}

func (s *classroomInvitationService) InviteByEmail(classroomID, inviterID, email string) (*dto.ClassroomInvitationResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if inviter is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, inviterID)
	if err != nil {
		return nil, err
	}
	if !isAllowed {
		return nil, apperror.ErrForbidden
	}

	// Get classroom info
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return nil, err
	}

	// Find user by email
	invitee, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	// Check if inviting self
	if invitee.ID.Hex() == inviterID {
		return nil, apperror.ErrCannotInviteSelf
	}

	// Check if user is already lecturer
	if classroom.LecturerID == invitee.ID {
		return nil, apperror.ErrAlreadyLecturer
	}

	// Check if user is already co-lecturer
	for _, coLecturerID := range classroom.CoLecturerIDs {
		if coLecturerID == invitee.ID {
			return nil, apperror.ErrAlreadyCoLecturer
		}
	}

	// Check if user is already a student
	for _, studentID := range classroom.StudentIDs {
		if studentID == invitee.ID {
			return nil, apperror.ErrAlreadyInClassroom
		}
	}

	// Check if invitation already exists
	existingInvitation, err := s.invitationRepo.GetByInviteeAndClassroom(ctx, invitee.ID.Hex(), classroomID)
	if err != nil {
		return nil, err
	}
	if existingInvitation != nil {
		return nil, apperror.ErrInvitationAlreadyExists
	}

	// Get inviter info
	inviter, err := s.userRepo.GetByID(ctx, inviterID)
	if err != nil {
		return nil, err
	}

	// Create invitation
	invitation := &model.ClassroomInvitation{
		ClassroomID: classroom.ID,
		InviterID:   inviter.ID,
		InviteeID:   invitee.ID,
	}

	createdInvitation, err := s.invitationRepo.Create(ctx, invitation)
	if err != nil {
		return nil, err
	}

	// Publish event for notification
	s.eventBus.Publish(bus.ClassroomInvitationEvent{
		InvitationID:  createdInvitation.ID.Hex(),
		ClassroomID:   classroom.ID.Hex(),
		ClassroomName: classroom.Name,
		InviterID:     inviter.ID.Hex(),
		InviterName:   inviter.FullName,
		InviteeID:     invitee.ID.Hex(),
	})

	return &dto.ClassroomInvitationResponse{
		ID:            createdInvitation.ID.Hex(),
		ClassroomID:   classroom.ID.Hex(),
		ClassroomName: classroom.Name,
		InviterID:     inviter.ID.Hex(),
		InviterName:   inviter.FullName,
		InviteeID:     invitee.ID.Hex(),
		InviteeName:   invitee.FullName,
		InviteeEmail:  invitee.Email,
		Status:        string(createdInvitation.Status),
		CreatedAt:     createdInvitation.CreatedAt,
	}, nil
}

func (s *classroomInvitationService) AcceptInvitation(invitationID, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperror.ErrUserNotFound
	}

	// Get invitation
	invitation, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}

	// Check if user is the invitee
	if invitation.InviteeID.Hex() != userID {
		return apperror.ErrForbidden
	}

	// Check if invitation is still pending
	if invitation.Status != model.InvitationPending {
		return apperror.ErrInvitationAlreadyProcessed
	}

	// Add user as co-lecturer using user ID only
	err = s.classroomRepo.AddCoLecturer(ctx, invitation.ClassroomID.Hex(), userID)
	if err != nil {
		return err
	}

	classroom, err := s.classroomRepo.GetByID(ctx, invitation.ClassroomID.Hex())
	if err != nil {
		return err
	}

	err = s.channelRepo.AddMember(ctx, classroom.GeneralChannelID.Hex(), user)
	if err != nil {
		return err
	}

	// Update invitation status
	err = s.invitationRepo.UpdateStatus(ctx, invitationID, model.InvitationAccepted)
	if err != nil {
		return err
	}

	return nil
}

func (s *classroomInvitationService) RejectInvitation(invitationID, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get invitation
	invitation, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}

	// Check if user is the invitee
	if invitation.InviteeID.Hex() != userID {
		return apperror.ErrForbidden
	}

	// Check if invitation is still pending
	if invitation.Status != model.InvitationPending {
		return apperror.ErrInvitationAlreadyProcessed
	}

	// Update invitation status
	err = s.invitationRepo.UpdateStatus(ctx, invitationID, model.InvitationRejected)
	if err != nil {
		return err
	}

	return nil
}

func (s *classroomInvitationService) CancelInvitation(invitationID, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get invitation
	invitation, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}

	// Check if user is lecturer or co-lecturer of the classroom
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, invitation.ClassroomID.Hex(), userID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	// Check if invitation is still pending
	if invitation.Status != model.InvitationPending {
		return apperror.ErrInvitationAlreadyProcessed
	}

	// Delete invitation
	return s.invitationRepo.Delete(ctx, invitationID)
}

func (s *classroomInvitationService) GetPendingInvitationsForUser(userID string, page, pageSize int) ([]dto.ClassroomInvitationResponse, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	invitations, total, err := s.invitationRepo.GetPendingByInvitee(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.ClassroomInvitationResponse, 0, len(invitations))
	for _, inv := range invitations {
		// Get classroom info
		classroom, err := s.classroomRepo.GetByID(ctx, inv.ClassroomID.Hex())
		if err != nil {
			continue
		}

		// Get inviter info
		inviter, err := s.userRepo.GetByID(ctx, inv.InviterID.Hex())
		if err != nil {
			continue
		}

		// Get invitee info
		invitee, err := s.userRepo.GetByID(ctx, inv.InviteeID.Hex())
		if err != nil {
			continue
		}

		responses = append(responses, dto.ClassroomInvitationResponse{
			ID:            inv.ID.Hex(),
			ClassroomID:   classroom.ID.Hex(),
			ClassroomName: classroom.Name,
			InviterID:     inviter.ID.Hex(),
			InviterName:   inviter.FullName,
			InviteeID:     invitee.ID.Hex(),
			InviteeName:   invitee.FullName,
			InviteeEmail:  invitee.Email,
			Status:        string(inv.Status),
			CreatedAt:     inv.CreatedAt,
		})
	}

	return responses, total, nil
}

func (s *classroomInvitationService) GetPendingInvitationsForClassroom(classroomID, userID string, page, pageSize int) ([]dto.ClassroomInvitationResponse, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, userID)
	if err != nil {
		return nil, 0, err
	}
	if !isAllowed {
		return nil, 0, apperror.ErrForbidden
	}

	// Get classroom info
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return nil, 0, err
	}

	invitations, total, err := s.invitationRepo.GetPendingByClassroom(ctx, classroomID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.ClassroomInvitationResponse, 0, len(invitations))
	for _, inv := range invitations {
		// Get inviter info
		inviter, err := s.userRepo.GetByID(ctx, inv.InviterID.Hex())
		if err != nil {
			continue
		}

		// Get invitee info
		invitee, err := s.userRepo.GetByID(ctx, inv.InviteeID.Hex())
		if err != nil {
			continue
		}

		responses = append(responses, dto.ClassroomInvitationResponse{
			ID:            inv.ID.Hex(),
			ClassroomID:   classroom.ID.Hex(),
			ClassroomName: classroom.Name,
			InviterID:     inviter.ID.Hex(),
			InviterName:   inviter.FullName,
			InviteeID:     invitee.ID.Hex(),
			InviteeName:   invitee.FullName,
			InviteeEmail:  invitee.Email,
			Status:        string(inv.Status),
			CreatedAt:     inv.CreatedAt,
		})
	}

	return responses, total, nil
}
