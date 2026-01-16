package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	RecipientID primitive.ObjectID     `bson:"recipient_id" json:"recipient_id"`
	ActorID     primitive.ObjectID     `bson:"actor_id,omitempty" json:"actor_id,omitempty"`
	Type        NotificationType       `bson:"type,omitempty" json:"type,omitempty"`
	Message     string                 `bson:"message,omitempty" json:"message,omitempty"`
	Link        string                 `bson:"link,omitempty" json:"link,omitempty"`
	IsRead      bool                   `bson:"is_read,omitempty" json:"is_read,omitempty"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt   time.Time              `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

type NotificationType string

const (
	NotificationTypeGroupInvite     NotificationType = "group_invite"
	NotificationTypeComment         NotificationType = "comment"
	NotificationTypeLike            NotificationType = "like"
	NotificationTypeFollow          NotificationType = "follow"
	NotificationTypeMention         NotificationType = "mention"
	NotificationTypeNewMessage      NotificationType = "new_message"
	NotificationTypeGroupInvitation NotificationType = "group_invitation"
	NotificationTypeSystem          NotificationType = "system"

	NotificationTypeProjectRegistrationDeadline NotificationType = "project_registration_deadline"
	NotificationTypeProjectRegistrationExpired  NotificationType = "project_registration_expired"
	NotificationTypeProjectRegistrationOpened   NotificationType = "project_registration_opened"

	NotificationTypeReportGraded    NotificationType = "report_graded"
	NotificationTypeReportExpired   NotificationType = "report_expired"
	NotificationTypeReportSubmitted NotificationType = "report_submitted"
	NotificationTypeReportDeadline  NotificationType = "report_deadline"
	NotificationTypeReportOpened    NotificationType = "report_opened"

	NotificationTypeClassUpdated NotificationType = "class_updated"
)
