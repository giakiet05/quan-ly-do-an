package dto

import (
	"time"

	"github.com/giakiet05/quan-ly-do-an/backend/internal/model"
)

// ==================== REQUEST DTOs ====================

// InviteRequest is used for manual bulk invite
type InviteRequest struct {
	Emails []string `json:"emails" binding:"required,min=1,dive,email"`
}

// AcceptInvitationRequest accepts an invitation by ID (no need for request body, use path param)
// RejectInvitationRequest rejects an invitation by ID (no need for request body, use path param)

// ==================== RESPONSE DTOs ====================

// InviteResult contains the result of invitation operations
type InviteResult struct {
	Invited        []string       `json:"invited"`
	AlreadyMembers []string       `json:"already_members"`
	Failed         []InviteError  `json:"failed"`
	TotalProcessed int            `json:"total_processed"`
}

// InviteError represents an error when processing an invitation
type InviteError struct {
	Row   int    `json:"row,omitempty"`   // For Excel import
	Email string `json:"email"`
	Error string `json:"error"`
}

// Removed InvitationInfoResponse - no longer need public endpoint for token-based preview

// InvitationResponse represents a classroom invitation
type InvitationResponse struct {
	ID            string                   `json:"id"`
	Email         string                   `json:"email"`
	ClassroomID   string                   `json:"classroom_id"`
	ClassroomName string                   `json:"classroom_name,omitempty"`
	InvitedBy     string                   `json:"invited_by"`
	Status        model.InvitationStatus   `json:"status"`
	ExpiresAt     time.Time                `json:"expires_at"`
	CreatedAt     time.Time                `json:"created_at"`
}

// AcceptInvitationResponse contains result of accepting an invitation
type AcceptInvitationResponse struct {
	ClassroomID   string `json:"classroom_id"`
	ClassroomName string `json:"classroom_name"`
}

// ==================== CONVERTERS ====================

func FromInvitation(inv *model.ClassroomInvitation) InvitationResponse {
	if inv == nil {
		return InvitationResponse{}
	}

	return InvitationResponse{
		ID:          inv.ID.Hex(),
		Email:       inv.Email,
		ClassroomID: inv.ClassroomID.Hex(),
		InvitedBy:   inv.InvitedBy.Hex(),
		Status:      inv.Status,
		ExpiresAt:   inv.ExpiresAt,
		CreatedAt:   inv.CreatedAt,
	}
}

func FromInvitations(invs []model.ClassroomInvitation) []InvitationResponse {
	responses := make([]InvitationResponse, 0, len(invs))
	for _, inv := range invs {
		responses = append(responses, FromInvitation(&inv))
	}
	return responses
}
