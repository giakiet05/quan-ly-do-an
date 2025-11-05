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
	Lecturer         UserInfo           `bson:"lecturer" json:"lecturer"`
	Students         []UserInfo         `bson:"students" json:"students"`
	Projects         []Project          `bson:"projects" json:"projects"`
	Setting          ClassroomSetting   `bson:"setting" json:"setting"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}

type Project struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	MinMember   int                `bson:"min_member" json:"min_member"`
	MaxMember   int                `bson:"max_member" json:"max_member"`
}

type ClassroomSetting struct {
	CanStudentMessage bool `bson:"can_student_message" json:"can_student_message"`
}
