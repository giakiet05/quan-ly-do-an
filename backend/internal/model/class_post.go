package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassPost struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClassroomID primitive.ObjectID `bson:"classroom_id" json:"classroom_id"`
	Author      UserInfo           `bson:"author" json:"author"`
	Title       string             `bson:"title" json:"title"`
	Content     string             `bson:"content" json:"content"`
	Attachments []Attachment       `bson:"attachments,omitempty" json:"attachments,omitempty"`
	IsPinned    bool               `bson:"is_pinned" json:"is_pinned"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Attachment struct {
	FileName string `bson:"file_name" json:"file_name"`
	FileURL  string `bson:"file_url" json:"file_url"`
	FileSize int64  `bson:"file_size" json:"file_size"` // bytes
	MimeType string `bson:"mime_type" json:"mime_type"`
}
