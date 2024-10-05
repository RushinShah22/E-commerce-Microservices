package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName       string             `bson:"firstName,omitempty" json:"firstName"`
	LastName        string             `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email"`
	Role            string             `bson:"role,omitempty" json:"role,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	ConfirmPassword string             `bson:"confirmPassword,omitempty" json:"confirmPassword,omitempty"`
	CreatedAt       primitive.DateTime `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}
