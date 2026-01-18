package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID             primitive.ObjectID    `bson:"_id,omitempty"`
	ClassroomID    primitive.ObjectID    `bson:"classroom_id,omitempty"`
	ProjectID      primitive.ObjectID    `bson:"project_id,omitempty"`
	GroupChannelID *primitive.ObjectID   `bson:"group_channel_id,omitempty"`
	LeaderID       primitive.ObjectID    `bson:"leader_id,omitempty"`
	Members        []UserInfo            `bson:"members" json:"members"`
	Tasks          []Task                `bson:"tasks" json:"tasks"`
	TaskStatuses   []string              `bson:"task_statuses" json:"task_statuses"`
	Reports        []Report              `bson:"reports" json:"reports"`
	Setting        GroupSetting          `bson:"setting" json:"setting"`
	MinMember      int                   `bson:"min_member" json:"min_member"`
	MaxMember      int                   `bson:"max_member" json:"max_member"`
	JoinRequests   []JoinGroupRequest    `bson:"join_requests" json:"join_requests"`
	JoinInvitation []JoinGroupInvitation `bson:"join_invitations" json:"join_invitations"`
}

type Task struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string               `bson:"title" json:"title"`
	Details     *string              `bson:"details,omitempty" json:"details,omitempty"`
	AssignToIDs []primitive.ObjectID `bson:"assign_to_ids" json:"assign_to_ids"`
	DueDate     time.Time            `bson:"due_date" json:"due_date"`
	Status      string               `bson:"status" json:"status"`
}

type Report struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	ReportPeriodID primitive.ObjectID `bson:"report_period_id" json:"report_period_id"`
	Title          string             `bson:"title" json:"title"`
	Content        string             `bson:"content" json:"content"`
	Files          []File             `bson:"files" json:"files"`
	Feedback       ReportFeedback     `bson:"feedback" json:"feedback"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type ReportFeedback struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Content     string             `bson:"content" json:"content"`
	Grade       string             `bson:"grade" json:"grade"`
	LecturerID  primitive.ObjectID `bson:"lecturer_id" json:"lecturer_id"`
	CommentedAt time.Time          `bson:"commented_at" json:"commented_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type GroupSetting struct {
	AllowJoinRequest bool `bson:"allow_join_request" json:"allow_join_request"`
}
