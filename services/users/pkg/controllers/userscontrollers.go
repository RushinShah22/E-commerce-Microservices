package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	database "github.com/RushinShah22/e-commerce-micro/services/users/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/users/pkg/models"
	"github.com/RushinShah22/e-commerce-micro/services/users/pkg/producers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllUser(w http.ResponseWriter, r *http.Request) {

	cursor, err := database.User.UserColl.Find(r.Context(), bson.M{})

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

	json.NewEncoder(w).Encode(users)

}

func GetAUser(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "id")

	id, _ := primitive.ObjectIDFromHex(_id)

	var user model.User

	database.User.UserColl.FindOne(r.Context(), bson.D{{Key: "_id", Value: id}}).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func AddAUser(w http.ResponseWriter, r *http.Request) {

	var user model.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		panic(err)
	}

	newUser, err := database.User.UserColl.InsertOne(r.Context(), user)

	if err != nil {
		panic(err)
	}
	user.ID = newUser.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	// Producing user
	go producers.ProduceMessage(user, producers.CREATED)

}
