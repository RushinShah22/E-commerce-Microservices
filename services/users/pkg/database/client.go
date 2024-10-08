package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserInterface struct {
	Client   *mongo.Client
	DB       *mongo.Database
	UserColl *mongo.Collection
}

var User UserInterface

func ConnToDB(uri string) {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// storing necessary data on the global var.
	User.Client = client
	User.DB = User.Client.Database("e-commerce")
	User.UserColl = User.DB.Collection("users")

}
