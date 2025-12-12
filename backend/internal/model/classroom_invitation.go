package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InvitationStatus defines the status of a classroom invitation.
type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusRejected InvitationStatus = "rejected"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusCanceled InvitationStatus = "canceled"
)

// ClassroomInvitation represents an invitation to join a classroom.
type ClassroomInvitation struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Email       string              `bson:"email" json:"email"`
	ClassroomID primitive.ObjectID  `bson:"classroom_id" json:"classroom_id"`
	InvitedBy   primitive.ObjectID  `bson:"invited_by" json:"invited_by"`
	InvitedTo   *primitive.ObjectID `bson:"invited_to,omitempty" json:"invited_to,omitempty"` // UserID if user exists
	Status      InvitationStatus    `bson:"status" json:"status"`
	ExpiresAt   time.Time           `bson:"expires_at" json:"expires_at"`
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
