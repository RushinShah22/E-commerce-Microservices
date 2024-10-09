package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name,omitempty" json:"name,omitempty"`
	Quantity int                `bson:"quantity,omitempty" json:"quantity,omitempty"`
	SellerID primitive.ObjectID `bson:"sellerID,omitempty" json:"sellerID,omitempty"`
}
