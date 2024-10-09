package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/database"
	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/model"
	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/producers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func checkQuantity(ctx context.Context, productID primitive.ObjectID, quantity int) int {
	var product model.Catalog
	database.Order.CatalogColl.FindOne(ctx, bson.D{{Key: "productID", Value: productID}}).Decode(&product)

	if product.ID == primitive.NilObjectID {
		return -1
	}
	return product.Quantity - quantity
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder model.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	newQuantity := checkQuantity(r.Context(), newOrder.ProductID, newOrder.Quantity)
	if newQuantity < 0 {
		http.Error(w, "Required Quantity or Product Not available.", http.StatusBadRequest)
		log.Print("Tried to buy more products than available.")
		return
	}

	var newProduct model.Catalog = model.Catalog{
		Quantity: newQuantity,
	}

	res, err := database.Order.OrderColl.InsertOne(r.Context(), newOrder)

	if err != nil {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	proRes := database.Order.CatalogColl.FindOneAndUpdate(r.Context(), bson.D{{Key: "productID", Value: newOrder.ProductID}}, bson.M{"$set": bson.M{"quantity": newQuantity}}, options.FindOneAndUpdate().SetReturnDocument(options.After))
	proRes.Decode(&newProduct)

	newOrder.ID = res.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&newOrder)

	// Producing new orders
	go producers.ProduceMessage(newOrder, producers.CREATED)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	var orders []model.Order
	w.Header().Set("Content-Type", "application/json")

	cursor, err := database.Order.OrderColl.Find(r.Context(), bson.M{})

	if err != nil {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	err = cursor.All(r.Context(), &orders)

	if err != nil {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func GetAOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	database.Order.OrderColl.FindOne(r.Context(), bson.D{{Key: "_id", Value: id}}).Decode(&order)

	if order.ID == primitive.NilObjectID {
		http.Error(w, "Order was not found", http.StatusNotFound)
		log.Printf("Tried to get a product with invalid id.")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)

}
