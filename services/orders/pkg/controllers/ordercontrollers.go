package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/database"
	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func checkQuantity(ctx context.Context, productID primitive.ObjectID, quantity int) bool {
	var product model.Catalog
	database.Order.CatalogColl.FindOne(ctx, bson.D{{Key: "productID", Value: productID}}).Decode(&product)
	return product.Quantity >= quantity
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder model.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)

	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if !checkQuantity(r.Context(), newOrder.ProductID, newOrder.Quantity) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Required Quantity Not available."))
	}

	res, err := database.Order.OrderColl.InsertOne(r.Context(), newOrder)
	newOrder.ID = res.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&newOrder)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	var orders []model.Order
	cursor, err := database.Order.OrderColl.Find(r.Context(), bson.M{})
	if err != nil {
		panic(err)
	}
	err = cursor.All(r.Context(), &orders)

	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", orders)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
