package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductID   primitive.ObjectID `bson:"productID,omitempty" json:"productID,omitempty"`
	Quantity    int                `bson:"quantity,omitempty" json:"quantity,omitempty"`
	UserID      primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	DateOfOrder primitive.DateTime `bson:"dateOfOrder,omitempty" json:"dateOfOrder,omitempty"`
	Status      string             `bson:"status,omitempty" json:"status,omitempty"`
}
