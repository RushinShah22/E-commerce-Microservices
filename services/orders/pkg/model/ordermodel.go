package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	ProductID   primitive.ObjectID `bson:"productID" json:"productID"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	DateOfOrder primitive.DateTime `bson:"dateOfOrder" json:"dateOfOrder"`
	Status      string             `bson:"status" json:"status"`
}
