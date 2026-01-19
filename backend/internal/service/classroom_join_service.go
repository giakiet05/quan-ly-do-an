package service

import (
	"errors"
	"strings"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassroomJoinService interface {
	PreviewClassroom(code string) (*dto.ClassroomPreviewResponse, error)
	JoinClassroom(userID string, req dto.JoinClassroomRequest) (*dto.JoinClassroomResponse, error)
	GetPendingRequests(classroomID, userID string, page, pageSize int) ([]dto.JoinRequestResponse, int64, error)
	ApproveRequest(requestID, reviewerID string) error
	RejectRequest(requestID, reviewerID string) error
}

type classroomJoinService struct {
	joinRequestRepo repo.ClassroomJoinRequestRepo
	classroomRepo   repo.ClassroomRepo
	userRepo        repo.UserRepo
}

func NewClassroomJoinService(
	joinRequestRepo repo.ClassroomJoinRequestRepo,
	classroomRepo repo.ClassroomRepo,
	userRepo repo.UserRepo,
) ClassroomJoinService {
	return &classroomJoinService{
		joinRequestRepo: joinRequestRepo,
		classroomRepo:   classroomRepo,
		userRepo:        userRepo,
	}
}

// PreviewClassroom returns classroom information for preview before joining
func (s *classroomJoinService) PreviewClassroom(code string) (*dto.ClassroomPreviewResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	classroom, err := s.classroomRepo.GetByInvitationCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &dto.ClassroomPreviewResponse{
		ID:           classroom.ID.Hex(),
		Name:         classroom.Name,
		Lecturer:     classroom.Lecturer.FullName,
		StudentCount: len(classroom.Students),
		MaxStudents:  classroom.MaxStudents,
	}, nil
}

// JoinClassroom handles classroom join request with full validation
func (s *classroomJoinService) JoinClassroom(userID string, req dto.JoinClassroomRequest) (*dto.JoinClassroomResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get classroom by invitation code
	classroom, err := s.classroomRepo.GetByInvitationCode(ctx, req.InvitationCode)
	if err != nil {
		return nil, err
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	// Validate email domain if restriction is enabled
	if classroom.EnableEmailRestriction && len(classroom.AllowedEmailDomains) > 0 {
		validEmail := false
		for _, domain := range classroom.AllowedEmailDomains {
			if strings.HasSuffix(user.Email, domain) {
				validEmail = true
				break
			}
		}
		if !validEmail {
			return nil, apperror.ErrInvalidEmailDomain
		}
	}

	// Only validate student code if user has one (skip for lecturers/TAs)
	if user.StudentCode != nil && *user.StudentCode != "" {
		// Validate student code in whitelist if whitelist is enabled
		if classroom.EnableWhitelist && len(classroom.WhitelistStudentCode) > 0 {
			if !containsWhitelistEntry(classroom.WhitelistStudentCode, *user.StudentCode) {
				return nil, apperror.ErrStudentCodeNotInWhitelist
			}

			// Check if student code has already been claimed in whitelist
			isClaimed, err := s.classroomRepo.IsStudentCodeClaimedInWhitelist(ctx, classroom.ID.Hex(), *user.StudentCode)
			if err != nil {
				return nil, err
			}
			if isClaimed {
				return nil, apperror.ErrStudentCodeAlreadyUsed
			}
		}

		// Check student code not already used in classroom (fallback check)
		exists, err := s.classroomRepo.StudentCodeExistsInClassroom(ctx, classroom.ID.Hex(), *user.StudentCode)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperror.ErrStudentCodeAlreadyUsed
		}
	}

	// Check classroom capacity
	if classroom.MaxStudents > 0 && len(classroom.Students) >= classroom.MaxStudents {
		return nil, apperror.ErrClassroomFull
	}

	// Check if user already in classroom
	isInClassroom, err := s.classroomRepo.IsStudentInClassroom(ctx, classroom.ID.Hex(), userID)
	if err != nil {
		return nil, err
	}
	if isInClassroom {
		return nil, apperror.ErrAlreadyInClassroom
	}

	// Check if user already has a pending/approved request
	existingRequest, err := s.joinRequestRepo.GetByUserAndClassroom(ctx, userID, classroom.ID.Hex())
	if err != nil {
		return nil, err
	}
	if existingRequest != nil {
		if existingRequest.Status == model.JoinRequestPending {
			return &dto.JoinClassroomResponse{
				ClassroomID: classroom.ID.Hex(),
				Status:      "pending",
				Message:     "You already have a pending join request for this classroom",
			}, nil
		}
		if existingRequest.Status == model.JoinRequestApproved {
			return nil, apperror.ErrAlreadyInClassroom
		}
	}

	// Auto-approve or create pending request
	if classroom.AutoApprove {
		// Add student directly
		studentInfo := model.UserInfo{
			ID:          user.ID,
			FullName:    user.FullName,
			Avatar:      user.Avatar,
			StudentCode: user.StudentCode,
		}

		err = s.classroomRepo.AddStudent(ctx, classroom.ID.Hex(), studentInfo)
		if err != nil {
			return nil, err
		}

		// Mark whitelist entry as joined if whitelist is enabled and user has student code
		if classroom.EnableWhitelist && user.StudentCode != nil && *user.StudentCode != "" {
			_ = s.classroomRepo.MarkWhitelistEntryAsJoined(ctx, classroom.ID.Hex(), *user.StudentCode, user.ID)
		}

		return &dto.JoinClassroomResponse{
			ClassroomID: classroom.ID.Hex(),
			Status:      "approved",
			Message:     "Successfully joined classroom",
		}, nil
	} else {
		// Create pending join request
		joinRequest := &model.ClassroomJoinRequest{
			ClassroomID:    classroom.ID,
			InvitationCode: req.InvitationCode,
			UserID:         user.ID,
			Status:         model.JoinRequestPending,
			CreatedAt:      time.Now(),
		}

		err = s.joinRequestRepo.Create(ctx, joinRequest)
		if err != nil {
			return nil, err
		}

		return &dto.JoinClassroomResponse{
			ClassroomID: classroom.ID.Hex(),
			Status:      "pending",
			Message:     "Join request submitted, waiting for lecturer approval",
		}, nil
	}
}

// GetPendingRequests returns pending join requests for a classroom (lecturer or co-lecturer)
func (s *classroomJoinService) GetPendingRequests(classroomID, userID string, page, pageSize int) ([]dto.JoinRequestResponse, int64, error) {
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

	// Get pending requests
	requests, total, err := s.joinRequestRepo.GetPendingByClassroom(ctx, classroomID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response DTOs with user info
	responses := make([]dto.JoinRequestResponse, 0, len(requests))
	for _, req := range requests {
		// Query user info
		user, err := s.userRepo.GetByID(ctx, req.UserID.Hex())
		if err != nil {
			// Skip if user not found (shouldn't happen)
			continue
		}

		studentCode := ""
		if user.StudentCode != nil {
			studentCode = *user.StudentCode
		}

		responses = append(responses, dto.JoinRequestResponse{
			ID:          req.ID.Hex(),
			StudentCode: studentCode,
			FullName:    user.FullName,
			Email:       user.Email,
			Status:      string(req.Status),
			CreatedAt:   req.CreatedAt,
		})
	}

	return responses, total, nil
}

// ApproveRequest approves a join request and adds student to classroom
func (s *classroomJoinService) ApproveRequest(requestID, reviewerID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get join request
	request, err := s.joinRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	// Check if already processed
	if request.Status != model.JoinRequestPending {
		return apperror.ErrJoinRequestAlreadyProcessed
	}

	// Get classroom
	classroom, err := s.classroomRepo.GetByID(ctx, request.ClassroomID.Hex())
	if err != nil {
		return err
	}

	// Verify reviewer is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroom.ID.Hex(), reviewerID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	reviewerOID, _ := primitive.ObjectIDFromHex(reviewerID)

	// Get user info
	user, err := s.userRepo.GetByID(ctx, request.UserID.Hex())
	if err != nil {
		return err
	}

	// Check student code not already used (only if user has student code)
	if user.StudentCode != nil && *user.StudentCode != "" {
		exists, err := s.classroomRepo.StudentCodeExistsInClassroom(ctx, classroom.ID.Hex(), *user.StudentCode)
		if err != nil {
			return err
		}
		if exists {
			return apperror.ErrStudentCodeAlreadyUsed
		}
	}

	// Check classroom capacity
	if classroom.MaxStudents > 0 && len(classroom.Students) >= classroom.MaxStudents {
		return apperror.ErrClassroomFull
	}

	// Add student to classroom
	studentInfo := model.UserInfo{
		ID:          user.ID,
		FullName:    user.FullName,
		Avatar:      user.Avatar,
		StudentCode: user.StudentCode,
	}

	err = s.classroomRepo.AddStudent(ctx, classroom.ID.Hex(), studentInfo)
	if err != nil {
		return err
	}

	// Mark whitelist entry as joined if whitelist is enabled and user has student code
	if classroom.EnableWhitelist && user.StudentCode != nil && *user.StudentCode != "" {
		_ = s.classroomRepo.MarkWhitelistEntryAsJoined(ctx, classroom.ID.Hex(), *user.StudentCode, user.ID)
	}

	// Update request status
	err = s.joinRequestRepo.UpdateStatus(ctx, requestID, model.JoinRequestApproved, &reviewerOID)
	if err != nil {
		return err
	}

	// TODO: Send notification to user

	return nil
}

// RejectRequest rejects a join request
func (s *classroomJoinService) RejectRequest(requestID, reviewerID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get join request
	request, err := s.joinRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	// Check if already processed
	if request.Status != model.JoinRequestPending {
		return apperror.ErrJoinRequestAlreadyProcessed
	}

	// Get classroom
	classroom, err := s.classroomRepo.GetByID(ctx, request.ClassroomID.Hex())
	if err != nil {
		return err
	}

	// Verify reviewer is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroom.ID.Hex(), reviewerID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	reviewerOID, _ := primitive.ObjectIDFromHex(reviewerID)

	// Update request status
	err = s.joinRequestRepo.UpdateStatus(ctx, requestID, model.JoinRequestRejected, &reviewerOID)
	if err != nil {
		return err
	}

	// TODO: Send notification to user

	return nil
}

// Helper functions

func containsWhitelistEntry(entries []model.WhitelistEntry, studentCode string) bool {
	for _, entry := range entries {
		if entry.StudentCode == studentCode {
			return true
		}
	}
	return false
}
