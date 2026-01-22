package service

import (
	"context"
	"errors"
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/apperror"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/dto"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/repo"
	"github.com/giakiet05/quan-ly-do-an/backend/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassroomService interface {
	CreateClassroom(req dto.CreateClassroomRequest, lecturerID string) (*dto.ClassroomResponse, error)
	GetClassroomByID(classroomID string) (*dto.ClassroomResponse, error)
	GetClassroomsByLecturer(lecturerID string, page, pageSize int) ([]dto.ClassroomResponse, int64, error)
	GetClassroomsByStudent(studentID string, page, pageSize int) ([]dto.ClassroomResponse, int64, error)
	UpdateClassroom(classroomID, lecturerID string, req dto.UpdateClassroomRequest) error
	UpdateClassroomStatus(classroomID, lecturerID string, status string) error
	DeleteClassroom(classroomID, lecturerID string) error
	RemoveStudentFromClassroom(classroomID, lecturerID, studentID string) error
	RemoveCoLecturerFromClassroom(classroomID, lecturerID, coLecturerID string) error
	LeaveClassroom(classroomID, userID string) error

	// Whitelist management
	UploadWhitelistStudentCode(classroomID, lecturerID string, studentCodes []string) error
	UpdateWhitelistStudentCode(classroomID, lecturerID string, addCodes, removeCodes []string) error
	ClearWhitelistStudentCode(classroomID, lecturerID string) error
	GetWhitelistStudentCode(classroomID, lecturerID string) ([]model.WhitelistEntry, error)

	RegenerateInvitationCode(classroomID, lecturerID string) (string, error)
}

type classroomService struct {
	classroomRepo repo.ClassroomRepo
	userRepo      repo.UserRepo
	channelRepo   repo.ChannelRepo
}

func NewClassroomService(
	classroomRepo repo.ClassroomRepo,
	userRepo repo.UserRepo,
	channelRepo repo.ChannelRepo,
) ClassroomService {
	return &classroomService{
		classroomRepo: classroomRepo,
		userRepo:      userRepo,
		channelRepo:   channelRepo,
	}
}

// CreateClassroom creates a new classroom
func (s *classroomService) CreateClassroom(req dto.CreateClassroomRequest, lecturerID string) (*dto.ClassroomResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get lecturer info
	lecturer, err := s.userRepo.GetByID(ctx, lecturerID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	// Generate invitation code (simple random code without university prefix)
	invitationCode := util.GenerateInvitationCode("")

	// Set default values
	maxStudents := req.MaxStudents
	if maxStudents == 0 {
		maxStudents = 100 // Default
	}

	// Create classroom model
	classroom := &model.Classroom{
		Name:                   req.Name,
		Description:            req.Description,
		Avatar:                 req.Avatar,
		Semester:               req.Semester,
		Year:                   req.Year,
		Status:                 model.ClassroomActive,
		LecturerID:             lecturer.ID,
		CoLecturerIDs:          []primitive.ObjectID{},
		StudentIDs:             []primitive.ObjectID{},
		ProjectRounds:          []model.ProjectRound{},
		InvitationCode:         invitationCode,
		MaxStudents:            maxStudents,
		AutoApprove:            req.AutoApprove,
		AllowedEmailDomains:    req.AllowedEmailDomains,
		EnableWhitelist:        req.EnableWhitelist,
		EnableEmailRestriction: req.EnableEmailRestriction,
		WhitelistStudentCode:   []model.WhitelistEntry{},
		CanStudentDeleteGroup:  false,
		CreatedAt:              time.Now(),
	}

	adminIDs := make([]string, 0, 1+len(classroom.CoLecturerIDs))
	adminIDs = append(adminIDs, classroom.LecturerID.Hex())
	for _, id := range classroom.CoLecturerIDs {
		adminIDs = append(adminIDs, id.Hex())
	}

	members := []model.UserInfo{
		{
			ID:       lecturer.ID,
			FullName: lecturer.FullName,
			Email:    lecturer.Email,
			Avatar:   lecturer.Avatar,
		},
	}

	generalChannel, err := s.channelRepo.CreateClassroomChannel(ctx, adminIDs, members)
	if err != nil {
		return nil, err
	}
	classroom.GeneralChannelID = generalChannel.ID

	// Save to database
	createdClassroom, err := s.classroomRepo.Create(ctx, classroom)
	if err != nil {
		return nil, err
	}

	// Convert to DTO response (no co-lecturers or students yet)
	response := dto.FromClassroomWithUsers(createdClassroom, lecturer, []*model.User{}, []*model.User{})
	return &response, nil
}

// GetClassroomByID retrieves a classroom by ID
func (s *classroomService) GetClassroomByID(classroomID string) (*dto.ClassroomResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get classroom
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return nil, err
	}

	// Populate users
	lecturer, coLecturers, students, err := s.PopulateClassroomUsers(ctx, classroom)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	response := dto.FromClassroomWithUsers(classroom, lecturer, coLecturers, students)
	return &response, nil
}

// GetClassroomsByLecturer retrieves classrooms by lecturer ID
func (s *classroomService) GetClassroomsByLecturer(lecturerID string, page, pageSize int) ([]dto.ClassroomResponse, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get classrooms
	classrooms, total, err := s.classroomRepo.GetByLecturer(ctx, lecturerID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Populate users for each classroom and convert to DTO
	responses := make([]dto.ClassroomResponse, 0, len(classrooms))
	for _, classroom := range classrooms {
		lecturer, coLecturers, students, err := s.PopulateClassroomUsers(ctx, &classroom)
		if err != nil {
			return nil, 0, err
		}
		response := dto.FromClassroomWithUsers(&classroom, lecturer, coLecturers, students)
		responses = append(responses, response)
	}

	return responses, total, nil
}

func (s *classroomService) GetClassroomsByStudent(studentID string, page, pageSize int) ([]dto.ClassroomResponse, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get classrooms
	classrooms, total, err := s.classroomRepo.GetByStudent(ctx, studentID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Populate users for each classroom and convert to DTO
	responses := make([]dto.ClassroomResponse, 0, len(classrooms))
	for _, classroom := range classrooms {
		lecturer, coLecturers, students, err := s.PopulateClassroomUsers(ctx, &classroom)
		if err != nil {
			return nil, 0, err
		}
		response := dto.FromClassroomWithUsers(&classroom, lecturer, coLecturers, students)
		responses = append(responses, response)
	}

	return responses, total, nil
}

// UpdateClassroom updates classroom information (lecturer or co-lecturer)
func (s *classroomService) UpdateClassroom(classroomID, lecturerID string, req dto.UpdateClassroomRequest) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	// Get classroom
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return err
	}

	// Update fields (only update non-nil values)
	if req.Name != nil {
		classroom.Name = *req.Name
	}
	if req.Description != nil {
		classroom.Description = *req.Description
	}
	if req.Avatar != nil {
		classroom.Avatar = *req.Avatar
	}
	if req.Semester != nil {
		classroom.Semester = *req.Semester
	}
	if req.Year != nil {
		classroom.Year = *req.Year
	}
	if req.MaxStudents != nil {
		classroom.MaxStudents = *req.MaxStudents
	}
	if req.AutoApprove != nil {
		classroom.AutoApprove = *req.AutoApprove
	}
	if req.AllowedEmailDomains != nil {
		classroom.AllowedEmailDomains = req.AllowedEmailDomains
	}
	if req.EnableWhitelist != nil {
		classroom.EnableWhitelist = *req.EnableWhitelist
	}
	if req.EnableEmailRestriction != nil {
		classroom.EnableEmailRestriction = *req.EnableEmailRestriction
	}
	if req.CanStudentDeleteGroup != nil {
		classroom.CanStudentDeleteGroup = *req.CanStudentDeleteGroup
	}

	// Update classroom in database
	_, err = s.classroomRepo.Update(ctx, classroom)
	return err
}

// UpdateClassroomStatus updates classroom status (lecturer or co-lecturer)
func (s *classroomService) UpdateClassroomStatus(classroomID, lecturerID string, status string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	// Get classroom
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return err
	}

	// Update status
	classroom.Status = model.ClassroomStatus(status)

	// Update classroom in database
	_, err = s.classroomRepo.Update(ctx, classroom)
	return err
}

// DeleteClassroom deletes a classroom (lecturer only)
func (s *classroomService) DeleteClassroom(classroomID, lecturerID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer
	isLecturer, err := s.classroomRepo.IsLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	// Delete classroom
	return s.classroomRepo.Delete(ctx, classroomID)
}

// RemoveStudentFromClassroom removes a student from classroom (lecturer or co-lecturer)
func (s *classroomService) RemoveStudentFromClassroom(classroomID, userID, studentID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, userID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	// Get student info to retrieve student_code
	student, err := s.userRepo.GetByID(ctx, studentID)
	if err != nil {
		return err
	}

	// Get classroom to check whitelist settings
	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return err
	}

	// Remove student
	err = s.classroomRepo.RemoveStudent(ctx, classroomID, studentID)
	if err != nil {
		return err
	}

	err = s.channelRepo.RemoveMember(ctx, classroom.GeneralChannelID.Hex(), userID)
	if err != nil {
		return err
	}

	// Reset whitelist entry if student has student_code and whitelist is enabled
	if classroom.EnableWhitelist && student.StudentCode != nil && *student.StudentCode != "" {
		// Reset the whitelist entry (set joined_by and joined_at to null)
		_ = s.classroomRepo.ResetWhitelistEntry(ctx, classroomID, *student.StudentCode)
	}

	return nil
}

// RemoveCoLecturerFromClassroom removes a co-lecturer from classroom (lecturer only)
func (s *classroomService) RemoveCoLecturerFromClassroom(classroomID, lecturerID, coLecturerID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Only the main lecturer can remove co-lecturers
	isLecturer, err := s.classroomRepo.IsLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return err
	}
	if !isLecturer {
		return apperror.ErrForbidden
	}

	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return err
	}

	// Remove co-lecturer
	err = s.classroomRepo.RemoveCoLecturer(ctx, classroomID, coLecturerID)
	if err != nil {
		return err
	}

	err = s.channelRepo.RemoveMember(ctx, classroom.GeneralChannelID.Hex(), coLecturerID)
	if err != nil {
		return err
	}

	return nil
}

// LeaveClassroom allows a student or co-lecturer to leave the classroom
func (s *classroomService) LeaveClassroom(classroomID, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return err
	}

	// Check if user is the main lecturer (cannot leave)
	isLecturer, err := s.classroomRepo.IsLecturer(ctx, classroomID, userID)
	if err != nil {
		return err
	}
	if isLecturer {
		return apperror.ErrForbidden // Lecturer cannot leave their own classroom
	}

	// Check if user is co-lecturer
	isCoLecturer, err := s.classroomRepo.IsCoLecturer(ctx, classroomID, userID)
	if err != nil {
		return err
	}
	if isCoLecturer {
		err = s.channelRepo.RemoveMember(ctx, classroom.GeneralChannelID.Hex(), userID)
		if err != nil {
			return err
		}

		return s.classroomRepo.RemoveCoLecturer(ctx, classroomID, userID)
	}

	// Check if user is student
	isStudent, err := s.classroomRepo.IsStudentInClassroom(ctx, classroomID, userID)
	if err != nil {
		return err
	}
	if isStudent {
		// Get student info to retrieve student_code
		student, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		// Remove student
		err = s.classroomRepo.RemoveStudent(ctx, classroomID, userID)
		if err != nil {
			return err
		}

		// Reset whitelist entry if student has student_code and whitelist is enabled
		if classroom.EnableWhitelist && student.StudentCode != nil && *student.StudentCode != "" {
			_ = s.classroomRepo.ResetWhitelistEntry(ctx, classroomID, *student.StudentCode)
		}

		err = s.channelRepo.RemoveMember(ctx, classroom.GeneralChannelID.Hex(), userID)
		if err != nil {
			return err
		}

		return nil
	}

	return apperror.ErrForbidden // User is not a member of this classroom
}

// UploadWhitelistStudentCode uploads student code whitelist (replaces all)
func (s *classroomService) UploadWhitelistStudentCode(classroomID, lecturerID string, studentCodes []string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	// Convert to WhitelistEntry structs
	entries := make([]model.WhitelistEntry, len(studentCodes))
	for i, code := range studentCodes {
		entries[i] = model.WhitelistEntry{StudentCode: code}
	}

	// Update whitelist (replace all)
	return s.classroomRepo.UpdateWhitelistStudentCode(ctx, classroomID, entries)
}

// UpdateWhitelistStudentCode adds or removes student codes from whitelist
func (s *classroomService) UpdateWhitelistStudentCode(classroomID, lecturerID string, addCodes, removeCodes []string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	// Add codes
	if len(addCodes) > 0 {
		if err := s.classroomRepo.AddToWhitelistStudentCode(ctx, classroomID, addCodes); err != nil {
			return err
		}
	}

	// Remove codes
	if len(removeCodes) > 0 {
		if err := s.classroomRepo.RemoveFromWhitelistStudentCode(ctx, classroomID, removeCodes); err != nil {
			return err
		}
	}

	return nil
}

// ClearWhitelistStudentCode clears all student codes from whitelist
func (s *classroomService) ClearWhitelistStudentCode(classroomID, lecturerID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return err
	}
	if !isAllowed {
		return apperror.ErrForbidden
	}

	return s.classroomRepo.ClearWhitelistStudentCode(ctx, classroomID)
}

// GetWhitelistStudentCode gets current whitelist
func (s *classroomService) GetWhitelistStudentCode(classroomID, lecturerID string) ([]model.WhitelistEntry, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return nil, err
	}
	if !isAllowed {
		return nil, apperror.ErrForbidden
	}

	classroom, err := s.classroomRepo.GetByID(ctx, classroomID)
	if err != nil {
		return nil, err
	}

	return classroom.WhitelistStudentCode, nil
}

// RegenerateInvitationCode generates a new invitation code
func (s *classroomService) RegenerateInvitationCode(classroomID, lecturerID string) (string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if user is lecturer or co-lecturer
	isAllowed, err := s.classroomRepo.IsLecturerOrCoLecturer(ctx, classroomID, lecturerID)
	if err != nil {
		return "", err
	}
	if !isAllowed {
		return "", apperror.ErrForbidden
	}

	// Generate new code
	newCode := util.GenerateInvitationCode("")

	// Update code
	err = s.classroomRepo.RegenerateInvitationCode(ctx, classroomID, newCode)
	if err != nil {
		return "", err
	}

	return newCode, nil
}

// Helper function to populate users for a classroom
func (s *classroomService) PopulateClassroomUsers(ctx context.Context, classroom *model.Classroom) (*model.User, []*model.User, []*model.User, error) {
	// Collect all user IDs
	userIDs := []string{classroom.LecturerID.Hex()}

	for _, coLecturerID := range classroom.CoLecturerIDs {
		userIDs = append(userIDs, coLecturerID.Hex())
	}

	for _, studentID := range classroom.StudentIDs {
		userIDs = append(userIDs, studentID.Hex())
	}

	// Batch query all users
	users, err := s.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create a map for quick lookup
	userMap := make(map[string]*model.User)
	for _, user := range users {
		if user != nil {
			userMap[user.ID.Hex()] = user
		}
	}

	// Get lecturer
	lecturer := userMap[classroom.LecturerID.Hex()]

	// Get co-lecturers
	coLecturers := make([]*model.User, 0, len(classroom.CoLecturerIDs))
	for _, coLecturerID := range classroom.CoLecturerIDs {
		if user, exists := userMap[coLecturerID.Hex()]; exists {
			coLecturers = append(coLecturers, user)
		}
	}

	// Get students
	students := make([]*model.User, 0, len(classroom.StudentIDs))
	for _, studentID := range classroom.StudentIDs {
		if user, exists := userMap[studentID.Hex()]; exists {
			students = append(students, user)
		}
	}

	return lecturer, coLecturers, students, nil
}
