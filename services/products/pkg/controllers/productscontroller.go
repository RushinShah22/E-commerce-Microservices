package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/RushinShah22/e-commerce-micro/services/products/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/products/pkg/models"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct model.Product

	json.NewDecoder(r.Body).Decode(&newProduct)
	insertedPro, err := database.Product.ProductColl.InsertOne(r.Context(), newProduct)

	if err != nil {
		panic(err)
	}
	newProduct.ID = insertedPro.InsertedID.(primitive.ObjectID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProduct)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []model.Product
	cursor, err := database.Product.ProductColl.Find(r.Context(), bson.D{})

	if err != nil {
		panic(err)
	}

	for cursor.Next(r.Context()) {
		var product model.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}

func GetAProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Something went Wrong.", http.StatusInternalServerError)
	}

	database.Product.ProductColl.FindOne(r.Context(), bson.D{{Key: "_id", Value: id}}).Decode(&product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
