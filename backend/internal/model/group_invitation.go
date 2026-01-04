package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JoinGroupInvitation struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	GroupID     primitive.ObjectID     `bson:"group_id" json:"group_id"`
	RecipientID primitive.ObjectID     `bson:"recipient_id" json:"recipient_id"`
	Status      JoinGroupRequestStatus `bson:"status" json:"status"`
	SentAt      time.Time              `bson:"sent_at" json:"sent_at"`
	RespondedAt *time.Time             `bson:"responded_at,omitempty" json:"responded_at,omitempty"`
}

type JoinGroupRequest struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID     `bson:"user_id" json:"user_id"`
	Status      JoinGroupRequestStatus `bson:"status" json:"status"`
	Message     string                 `bson:"message" json:"message"`
	RequestedAt time.Time              `bson:"requested_at" json:"requested_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}

type JoinGroupRequestStatus string

const (
	RequestPending  JoinGroupRequestStatus = "pending"
	RequestAccepted JoinGroupRequestStatus = "accepted"
	RequestRejected JoinGroupRequestStatus = "rejected"
)
