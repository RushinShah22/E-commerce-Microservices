package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RushinShah22/e-commerce-micro/services/products/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/products/pkg/models"
	"github.com/RushinShah22/e-commerce-micro/services/products/pkg/producers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This is helper function to fetch a particular seller
func GetSeller(ctx context.Context, id primitive.ObjectID) model.Seller {
	var seller model.Seller
	database.Product.SellerColl.FindOne(ctx, bson.D{{Key: "userID", Value: id}}).Decode(&seller)
	return seller
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct model.Product
	w.Header().Set("Content-Type", "application/json")

	json.NewDecoder(r.Body).Decode(&newProduct)
	seller := GetSeller(r.Context(), newProduct.SellerID)

	if seller.ID == primitive.NilObjectID {
		log.Println("Tried to create product with invalid seller.")
		http.Error(w, "No seller exists with the provided id.", http.StatusBadRequest)
		return
	}
	insertedPro, err := database.Product.ProductColl.InsertOne(r.Context(), newProduct)

	// inserting the new document id in the struct
	newProduct.ID = insertedPro.InsertedID.(primitive.ObjectID)
	if err != nil || newProduct.ID == primitive.NilObjectID {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)

	// Produce message in Go routine
	go producers.ProduceMessage(newProduct, producers.CREATED)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []model.Product
	cursor, err := database.Product.ProductColl.Find(r.Context(), bson.D{})

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	for cursor.Next(r.Context()) {
		var product model.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	if err != nil {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}

func GetAProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, "Something went Wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	database.Product.ProductColl.FindOne(r.Context(), bson.D{{Key: "_id", Value: id}}, options.FindOne()).Decode(&product)
	w.Header().Set("Content-Type", "application/json")
	if product.ID == primitive.NilObjectID {
		log.Printf("Made request with wrong product ID: %s\n", id)
		http.Error(w, "No such product exists with id: "+id.Hex(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// This function is used for updating various parameters of a product.
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var newDetails model.Product
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, "Something went Wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	json.NewDecoder(r.Body).Decode(&newDetails)

	var updatedProduct model.Product

	database.Product.ProductColl.FindOneAndUpdate(r.Context(), bson.D{{Key: "_id", Value: id}}, bson.D{{Key: "$set", Value: newDetails}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedProduct)

	if updatedProduct.ID == primitive.NilObjectID {
		log.Printf("Made request with wrong product ID: %s\n", id)
		http.Error(w, "No such product exists with id: "+id.Hex(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(updatedProduct)

	// Produce message in Go routine
	go producers.ProduceMessage(updatedProduct, producers.UPDATED)
}
