package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type University struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
	Code string             `bson:"code" json:"code"`
}
