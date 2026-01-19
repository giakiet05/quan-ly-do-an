package dto

import "time"

// ===== REQUEST DTOs =====

type InviteToClassroomRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ===== RESPONSE DTOs =====

type ClassroomInvitationResponse struct {
	ID            string    `json:"id"`
	ClassroomID   string    `json:"classroom_id"`
	ClassroomName string    `json:"classroom_name"`
	InviterID     string    `json:"inviter_id"`
	InviterName   string    `json:"inviter_name"`
	InviteeID     string    `json:"invitee_id"`
	InviteeName   string    `json:"invitee_name"`
	InviteeEmail  string    `json:"invitee_email"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
