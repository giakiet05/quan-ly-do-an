package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassroomStatus string

const (
	ClassroomActive   ClassroomStatus = "active"
	ClassroomInactive ClassroomStatus = "inactive"
	ClassroomArchived ClassroomStatus = "archived"
)

type Classroom struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Description      string             `bson:"description" json:"description"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	Semester         string             `bson:"semester" json:"semester"` // "Fall 2026", "Spring 2027"
	Year             int                `bson:"year" json:"year"`
	Status           ClassroomStatus    `bson:"status" json:"status"` // "active", "inactive", "archived"
	GeneralChannelID primitive.ObjectID `bson:"general_channel_id" json:"general_channel_id"`
	Lecturer         UserInfo           `bson:"lecturer" json:"lecturer"`
	Students         []UserInfo         `bson:"students" json:"students"`
	ProjectRounds    []ProjectRound     `bson:"project_rounds" json:"project_rounds"`
	
	// Settings
	InvitationCode        string     `bson:"invitation_code" json:"invitation_code"`
	WhitelistStudentCode  []string   `bson:"whitelist_student_code,omitempty" json:"whitelist_student_code,omitempty"`
	RequireEmailDomain    *string    `bson:"require_email_domain,omitempty" json:"require_email_domain,omitempty"`
	MaxStudents           int        `bson:"max_students" json:"max_students"`
	AutoApprove           bool       `bson:"auto_approve" json:"auto_approve"`
	CanStudentDeleteGroup bool       `bson:"can_student_delete_group" json:"can_student_delete_group"`
	
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
