package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JoinRequestStatus string

const (
	JoinRequestPending  JoinRequestStatus = "pending"
	JoinRequestApproved JoinRequestStatus = "approved"
	JoinRequestRejected JoinRequestStatus = "rejected"
)

type ClassroomJoinRequest struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ClassroomID    primitive.ObjectID  `bson:"classroom_id" json:"classroom_id"`
	InvitationCode string              `bson:"invitation_code" json:"invitation_code"`
	UserID         primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Status         JoinRequestStatus   `bson:"status" json:"status"`
	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
	ReviewedAt     *time.Time          `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`
	ReviewedBy     *primitive.ObjectID `bson:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`
}
