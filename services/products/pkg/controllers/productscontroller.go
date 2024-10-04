package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/RushinShah22/e-commerce-micro/services/products/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/products/pkg/models"
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
	json.NewEncoder(w).Encode(newProduct)
}
