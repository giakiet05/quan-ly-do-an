package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Classroom struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UniversityID     primitive.ObjectID `bson:"university_id" json:"university_id"`
	UniversityName   string             `bson:"university_name" json:"university_name"`
	Name             string             `bson:"name" json:"name"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	GeneralChannelID primitive.ObjectID `bson:"general_channel_id" json:"general_channel_id"`
	Lecturer         UserInfo           `bson:"lecturer" json:"lecturer"`
	Students         []UserInfo         `bson:"students" json:"students"`
	ProjectRounds    []ProjectRound     `bson:"project_rounds" json:"project_rounds"`
	Setting          ClassroomSetting   `bson:"setting" json:"setting"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}

type ClassroomSetting struct {
	CanStudentDeleteGroup bool `bson:"can_student_delete_group" json:"can_student_delete_group"`
}
