package dto

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

// ===== REQUEST DTOs =====

type CreateClassroomRequest struct {
	Name               string  `json:"name" binding:"required"`
	Description        string  `json:"description"`
	Avatar             string  `json:"avatar"`
	Semester           string  `json:"semester"`           // "Fall 2026", "Spring 2027"
	Year               int     `json:"year"`               // 2026, 2027
	MaxStudents        int     `json:"max_students"`       // Optional, defaults to 100
	AutoApprove        bool    `json:"auto_approve"`       // Optional, defaults to false
	RequireEmailDomain *string `json:"require_email_domain"` // Optional, e.g., "@hcmut.edu.vn"
}

type UpdateClassroomRequest struct {
	Name                  *string `json:"name"`
	Description           *string `json:"description"`
	Avatar                *string `json:"avatar"`
	Semester              *string `json:"semester"`
	Year                  *int    `json:"year"`
	MaxStudents           *int    `json:"max_students"`
	AutoApprove           *bool   `json:"auto_approve"`
	RequireEmailDomain    *string `json:"require_email_domain"`
	CanStudentDeleteGroup *bool   `json:"can_student_delete_group"`
}

type UpdateClassroomStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive archived"`
}

type UploadWhitelistStudentCodeRequest struct {
	StudentCodes []string `json:"student_codes" binding:"required,min=1"`
}

type UpdateWhitelistStudentCodeRequest struct {
	AddCodes    []string `json:"add_codes"`    // Student codes to add
	RemoveCodes []string `json:"remove_codes"` // Student codes to remove
}

// ===== RESPONSE DTOs =====

type ClassroomResponse struct {
	ID                   string                   `json:"id"`
	Name                 string                   `json:"name"`
	Description          string                   `json:"description"`
	Avatar               string                   `json:"avatar"`
	Semester             string                   `json:"semester"`
	Year                 int                      `json:"year"`
	Status               string                   `json:"status"`
	Lecturer             model.UserInfoResponse   `json:"lecturer"`
	Students             []model.UserInfoResponse `json:"students"`
	InvitationCode       string                   `json:"invitation_code"`
	MaxStudents          int                      `json:"max_students"`
	AutoApprove          bool                     `json:"auto_approve"`
	WhitelistStudentCode []string                 `json:"whitelist_student_code,omitempty"`
	CreatedAt            time.Time                `json:"created_at"`
}

type RegenerateCodeResponse struct {
	NewCode       string    `json:"new_code"`
	CodeExpiresAt time.Time `json:"code_expires_at"`
}

// ===== CONVERTERS =====

func FromClassroom(classroom *model.Classroom) ClassroomResponse {
	if classroom == nil {
		return ClassroomResponse{}
	}

	students := make([]model.UserInfoResponse, 0, len(classroom.Students))
	for _, student := range classroom.Students {
		students = append(students, model.UserInfoResponse{
			UserID:   student.ID.Hex(),
			FullName: student.FullName,
			Avatar:   student.Avatar,
		})
	}

	return ClassroomResponse{
		ID:                   classroom.ID.Hex(),
		Name:                 classroom.Name,
		Description:          classroom.Description,
		Avatar:               classroom.Avatar,
		Semester:             classroom.Semester,
		Year:                 classroom.Year,
		Status:               string(classroom.Status),
		Lecturer: model.UserInfoResponse{
			UserID:   classroom.Lecturer.ID.Hex(),
			FullName: classroom.Lecturer.FullName,
			Avatar:   classroom.Lecturer.Avatar,
		},
		Students:             students,
		InvitationCode:       classroom.InvitationCode,
		MaxStudents:          classroom.MaxStudents,
		AutoApprove:          classroom.AutoApprove,
		WhitelistStudentCode: classroom.WhitelistStudentCode,
		CreatedAt:            classroom.CreatedAt,
	}
}

func FromClassrooms(classrooms []model.Classroom) []ClassroomResponse {
	responses := make([]ClassroomResponse, 0, len(classrooms))
	for _, classroom := range classrooms {
		responses = append(responses, FromClassroom(&classroom))
	}
	return responses
}
