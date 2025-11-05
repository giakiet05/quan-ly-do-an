package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Membership struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id,omitempty"`
	ClassRoomID primitive.ObjectID `bson:"classroom_id,omitempty"`
}
