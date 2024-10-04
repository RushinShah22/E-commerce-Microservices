package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderInterface struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var Order OrderInterface

func ConnToDB(uri string) {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	Order.Client = client
	Order.DB = Order.Client.Database("e-commerce")
}
