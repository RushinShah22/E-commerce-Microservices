package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Seller struct {
	UserID    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstName,omitempty" json:"firstName,omitempty"`
	LastName  string             `bson:"lastName,omitempty"  json:"lastName,omitempty"`
}
