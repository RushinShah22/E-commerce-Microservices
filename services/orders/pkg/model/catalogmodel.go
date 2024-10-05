package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Catalog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductID primitive.ObjectID `bson:"productID,omitempty" json:"productID,omitempty"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	Quantity  int                `bson:"quantity,omitempty" json:"quantity,omitempty"`
}
