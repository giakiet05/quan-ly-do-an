package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Classroom struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UniversityID     primitive.ObjectID `bson:"university_id" json:"university_id"`
	UniversityName   string             `bson:"university_name" json:"university_name"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	GeneralChannelID primitive.ObjectID `bson:"general_channel_id" json:"general_channel_id"`
	Lecturers        []UserInfo         `bson:"lecturers" json:"lecturers"`
	Students         []UserInfo         `bson:"students" json:"students"`
	ProjectGroups    []ProjectGroup     `bson:"project_groups" json:"project_groups"`
	Setting          ClassroomSetting   `bson:"setting" json:"setting"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}

type ProjectGroup struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Projects    []Project          `bson:"projects" json:"projects"`
	IsDeleted   bool               `bson:"is_deleted" json:"is_deleted"`
	DeletedAt   time.Time          `bson:"deleted_at" json:"deleted_at"`
}

type Project struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	MinMember   int                `bson:"min_member" json:"min_member"`
	MaxMember   int                `bson:"max_member" json:"max_member"`
}

type ClassroomSetting struct {
	CanStudentMessage     bool `bson:"can_student_message" json:"can_student_message"`
	CanStudentDeleteGroup bool `bson:"can_student_delete_group" json:"can_student_delete_group"`
	CanGroupChangeProject bool `bson:"can_group_change_project" json:"can_group_change_project"`
}
