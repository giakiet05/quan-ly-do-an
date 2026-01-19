package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvitationStatus string

const (
	InvitationPending  InvitationStatus = "pending"
	InvitationAccepted InvitationStatus = "accepted"
	InvitationRejected InvitationStatus = "rejected"
)

type ClassroomInvitation struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClassroomID primitive.ObjectID `bson:"classroom_id" json:"classroom_id"`
	InviterID   primitive.ObjectID `bson:"inviter_id" json:"inviter_id"`
	InviteeID   primitive.ObjectID `bson:"invitee_id" json:"invitee_id"`
	Status      InvitationStatus   `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	RespondedAt *time.Time         `bson:"responded_at,omitempty" json:"responded_at,omitempty"`
}
