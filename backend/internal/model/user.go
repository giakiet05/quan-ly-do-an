package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthProvider defines the source of user authentication.
type AuthProvider string

const (
	ProviderLocal  AuthProvider = "local"  // Registered with email and password
	ProviderGoogle AuthProvider = "google" // Registered via Google OAuth
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password,omitempty" json:"-"`
	AuthProvider AuthProvider       `bson:"auth_provider" json:"auth_provider"`
	ProviderID   string             `bson:"provider_id,omitempty" json:"-"`
	Avatar       *Image             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	StudentCode  *string            `bson:"student_code,omitempty" json:"student_code,omitempty"`
	FullName     string             `bson:"full_name" json:"full_name"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // Soft delete by user
}
