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

type WhitelistEntry struct {
	StudentCode string              `bson:"student_code" json:"student_code"`
	JoinedBy    *primitive.ObjectID `bson:"joined_by,omitempty" json:"joined_by,omitempty"`
	JoinedAt    *time.Time          `bson:"joined_at,omitempty" json:"joined_at,omitempty"`
}

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
	CoLecturers      []UserInfo         `bson:"co_lecturers" json:"co_lecturers"`
	Students         []UserInfo         `bson:"students" json:"students"`
	ProjectRounds    []ProjectRound     `bson:"project_rounds" json:"project_rounds"`

	// Settings
	InvitationCode         string           `bson:"invitation_code" json:"invitation_code"`
	WhitelistStudentCode   []WhitelistEntry `bson:"whitelist_student_code,omitempty" json:"whitelist_student_code,omitempty"`
	AllowedEmailDomains    []string         `bson:"allowed_email_domains,omitempty" json:"allowed_email_domains,omitempty"`
	EnableWhitelist        bool             `bson:"enable_whitelist" json:"enable_whitelist"`
	EnableEmailRestriction bool             `bson:"enable_email_restriction" json:"enable_email_restriction"`
	MaxStudents            int              `bson:"max_students" json:"max_students"`
	AutoApprove            bool             `bson:"auto_approve" json:"auto_approve"`
	CanStudentDeleteGroup  bool             `bson:"can_student_delete_group" json:"can_student_delete_group"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
