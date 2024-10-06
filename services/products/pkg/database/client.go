package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductInterface struct {
	Client      *mongo.Client
	DB          *mongo.Database
	ProductColl *mongo.Collection
	SellerColl  *mongo.Collection
}

var Product ProductInterface

func ConnToDB(uri string) {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	Product.Client = client
	Product.DB = Product.Client.Database("e-commerce")
	Product.ProductColl = Product.DB.Collection("products")
	Product.SellerColl = Product.DB.Collection("sellers")
}
