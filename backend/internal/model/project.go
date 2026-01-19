package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectRound struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	StartDate     time.Time          `bson:"start_date" json:"start_date"`
	EndDate       time.Time          `bson:"end_date" json:"end_date"`
	Description   string             `bson:"description" json:"description"`
	ReportPeriods []ReportPeriod     `bson:"report_periods" json:"report_periods"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	IsDeleted     bool               `bson:"is_deleted" json:"is_deleted"`
}

type Project struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClassroomID    primitive.ObjectID `bson:"classroom_id" json:"classroom_id"`
	ProjectRoundID primitive.ObjectID `bson:"project_round_id" json:"project_round_id"`
	Title          string             `bson:"title" json:"title"`
	Amount         int                `bson:"amount" json:"amount"`
	Description    string             `bson:"description" json:"description"`
	MinMember      int                `bson:"min_member" json:"min_member"`
	MaxMember      int                `bson:"max_member" json:"max_member"`
	Status         ProjectStatus      `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
}

type ProjectStatus string

const (
	ProjectStatusPending   ProjectStatus = "pending"
	ProjectStatusApproved  ProjectStatus = "approved"
	ProjectStatusOngoing   ProjectStatus = "ongoing"
	ProjectStatusCompleted ProjectStatus = "completed"
)

type ReportPeriod struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	FileType    []string           `bson:"file_type" json:"file_type"`
	StartDate   time.Time          `bson:"start_date" json:"start_date"`
	EndDate     time.Time          `bson:"end_date" json:"end_date"`
}
