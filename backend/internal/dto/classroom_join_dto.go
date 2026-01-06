package dto

import "time"

// ===== REQUEST DTOs =====

type JoinClassroomRequest struct {
	InvitationCode string `json:"invitation_code" binding:"required"`
}

type UploadWhitelistRequest struct {
	StudentCodes []string `json:"student_codes" binding:"required,min=1"`
}

// ===== RESPONSE DTOs =====

type ClassroomPreviewResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Lecturer     string `json:"lecturer"`
	StudentCount int    `json:"student_count"`
	MaxStudents  int    `json:"max_students"`
}

type JoinClassroomResponse struct {
	ClassroomID string `json:"classroom_id"`
	Status      string `json:"status"` // "approved" or "pending"
	Message     string `json:"message"`
}

type JoinRequestResponse struct {
	ID          string    `json:"id"`
	StudentCode string    `json:"student_code"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
