package dto

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

// ===== REQUEST DTOs =====

type CreateClassroomRequest struct {
	Name                   string   `json:"name" binding:"required"`
	Description            string   `json:"description"`
	Avatar                 string   `json:"avatar"`
	Semester               string   `json:"semester"`                 // "Fall 2026", "Spring 2027"
	Year                   int      `json:"year"`                     // 2026, 2027
	MaxStudents            int      `json:"max_students"`             // Optional, defaults to 100
	AutoApprove            bool     `json:"auto_approve"`             // Optional, defaults to false
	AllowedEmailDomains    []string `json:"allowed_email_domains"`    // Optional, e.g., ["@hcmut.edu.vn", "@student.hcmut.edu.vn"]
	EnableWhitelist        bool     `json:"enable_whitelist"`         // true = MSSV must be in whitelist
	EnableEmailRestriction bool     `json:"enable_email_restriction"` // true = email must match allowed domains
}

type UpdateClassroomRequest struct {
	Name                   *string  `json:"name"`
	Description            *string  `json:"description"`
	Avatar                 *string  `json:"avatar"`
	Semester               *string  `json:"semester"`
	Year                   *int     `json:"year"`
	MaxStudents            *int     `json:"max_students"`
	AutoApprove            *bool    `json:"auto_approve"`
	AllowedEmailDomains    []string `json:"allowed_email_domains"`
	EnableWhitelist        *bool    `json:"enable_whitelist"`
	EnableEmailRestriction *bool    `json:"enable_email_restriction"`
	CanStudentDeleteGroup  *bool    `json:"can_student_delete_group"`
}

type UpdateClassroomStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive archived"`
}

type UploadWhitelistStudentCodeRequest struct {
	StudentCodes []string `json:"student_codes" binding:"required,min=1"`
}

type WhitelistEntryResponse struct {
	StudentCode string  `json:"student_code"`
	JoinedBy    *string `json:"joined_by,omitempty"`
	JoinedAt    *string `json:"joined_at,omitempty"`
}

type UpdateWhitelistStudentCodeRequest struct {
	AddCodes    []string `json:"add_codes"`    // Student codes to add
	RemoveCodes []string `json:"remove_codes"` // Student codes to remove
}

// ===== RESPONSE DTOs =====

type ClassroomResponse struct {
	ID                     string                   `json:"id"`
	Name                   string                   `json:"name"`
	Description            string                   `json:"description"`
	Avatar                 string                   `json:"avatar"`
	Semester               string                   `json:"semester"`
	Year                   int                      `json:"year"`
	Status                 string                   `json:"status"`
	Lecturer               model.UserInfoResponse   `json:"lecturer"`
	CoLecturers            []model.UserInfoResponse `json:"co_lecturers"`
	Students               []model.UserInfoResponse `json:"students"`
	GeneralChannelID       string                   `json:"general_channel_id"`
	InvitationCode         string                   `json:"invitation_code"`
	MaxStudents            int                      `json:"max_students"`
	AutoApprove            bool                     `json:"auto_approve"`
	WhitelistStudentCode   []WhitelistEntryResponse `json:"whitelist_student_code,omitempty"`
	AllowedEmailDomains    []string                 `json:"allowed_email_domains,omitempty"`
	EnableWhitelist        bool                     `json:"enable_whitelist"`
	EnableEmailRestriction bool                     `json:"enable_email_restriction"`
	CreatedAt              time.Time                `json:"created_at"`
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

	coLecturers := make([]model.UserInfoResponse, 0, len(classroom.CoLecturers))
	for _, coLecturer := range classroom.CoLecturers {
		coLecturers = append(coLecturers, model.UserInfoResponse{
			UserID:   coLecturer.ID.Hex(),
			FullName: coLecturer.FullName,
			Avatar:   coLecturer.Avatar,
		})
	}

	students := make([]model.UserInfoResponse, 0, len(classroom.Students))
	for _, student := range classroom.Students {
		students = append(students, model.UserInfoResponse{
			UserID:   student.ID.Hex(),
			FullName: student.FullName,
			Avatar:   student.Avatar,
		})
	}

	whitelist := make([]WhitelistEntryResponse, 0, len(classroom.WhitelistStudentCode))
	for _, entry := range classroom.WhitelistStudentCode {
		var joinedBy *string
		var joinedAt *string
		if entry.JoinedBy != nil {
			hex := entry.JoinedBy.Hex()
			joinedBy = &hex
		}
		if entry.JoinedAt != nil {
			formatted := entry.JoinedAt.Format(time.RFC3339)
			joinedAt = &formatted
		}
		whitelist = append(whitelist, WhitelistEntryResponse{
			StudentCode: entry.StudentCode,
			JoinedBy:    joinedBy,
			JoinedAt:    joinedAt,
		})
	}

	return ClassroomResponse{
		ID:          classroom.ID.Hex(),
		Name:        classroom.Name,
		Description: classroom.Description,
		Avatar:      classroom.Avatar,
		Semester:    classroom.Semester,
		Year:        classroom.Year,
		Status:      string(classroom.Status),
		Lecturer: model.UserInfoResponse{
			UserID:   classroom.Lecturer.ID.Hex(),
			FullName: classroom.Lecturer.FullName,
			Avatar:   classroom.Lecturer.Avatar,
		},
		CoLecturers:            coLecturers,
		Students:               students,
		GeneralChannelID:       classroom.GeneralChannelID.Hex(),
		InvitationCode:         classroom.InvitationCode,
		MaxStudents:            classroom.MaxStudents,
		AutoApprove:            classroom.AutoApprove,
		WhitelistStudentCode:   whitelist,
		AllowedEmailDomains:    classroom.AllowedEmailDomains,
		EnableWhitelist:        classroom.EnableWhitelist,
		EnableEmailRestriction: classroom.EnableEmailRestriction,
		CreatedAt:              classroom.CreatedAt,
	}
}

func FromClassrooms(classrooms []model.Classroom) []ClassroomResponse {
	responses := make([]ClassroomResponse, 0, len(classrooms))
	for _, classroom := range classrooms {
		responses = append(responses, FromClassroom(&classroom))
	}
	return responses
}
