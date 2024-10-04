package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	database "github.com/RushinShah22/e-commerce-micro/services/users/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/users/pkg/models"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	coll := database.User.DB.Collection("users")

	cursor, err := coll.Find(r.Context(), bson.M{})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	var users []model.User
	err = cursor.All(r.Context(), &users)

	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	usersList, err := json.Marshal(users)
	w.Write(usersList)

}

func GetAUser(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "id")

	coll := database.User.DB.Collection("users")
	id, _ := primitive.ObjectIDFromHex(_id)

	var user model.User
	coll.FindOne(r.Context(), bson.D{{Key: "_id", Value: id}}).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)

}

func AddAUser(w http.ResponseWriter, r *http.Request) {
	coll := database.User.DB.Collection("users")
	var user model.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		panic(err)
	}

	newUser, err := coll.InsertOne(r.Context(), user)

	if err != nil {
		panic(err)
	}
	user.ID = newUser.InsertedID.(primitive.ObjectID)
	userJson, err := json.MarshalIndent(user, " ", " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJson)

}
