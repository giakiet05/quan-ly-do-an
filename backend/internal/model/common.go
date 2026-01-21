package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Image defines the structure for a stored image.
type Image struct {
	URL        string    `bson:"url" json:"url"`
	PublicID   string    `bson:"public_id" json:"public_id"`
	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`
}

// Video defines the structure for a stored video.
type Video struct {
	URL        string    `bson:"url" json:"url"`
	PublicID   string    `bson:"public_id" json:"public_id"`
	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`
}

type File struct{}

type UserInfo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName    string             `bson:"full_name" json:"full_name"`
	Email       string             `bson:"email" json:"email"`
	Avatar      *Image             `bson:"avatar" json:"avatar"`
	StudentCode *string            `bson:"student_code,omitempty" json:"student_code,omitempty"`
}
