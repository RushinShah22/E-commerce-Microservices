package client

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserInstance mongo.Client

func ConnToDB(uri string) {

	UserInstance, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := UserInstance.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
