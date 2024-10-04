package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Catalog struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	ProductID   primitive.ObjectID `bson:"productID" json:"productID"`
	ProductName string             `bson:"productName" json:"productName"`
	Quantity    int                `bson:"quantity" json:"quantity"`
}
