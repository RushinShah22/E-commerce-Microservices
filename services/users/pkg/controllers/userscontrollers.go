package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	database "github.com/RushinShah22/e-commerce-micro/services/users/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/users/pkg/models"
	"github.com/RushinShah22/e-commerce-micro/services/users/pkg/producers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cursor, err := database.User.UserColl.Find(r.Context(), bson.M{}, options.Find().SetProjection(bson.D{{Key: "password", Value: 0}, {Key: "createdAt", Value: 0}}))

	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	var users []model.User
	err = cursor.All(r.Context(), &users)

	if err != nil {
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(users)

}

func GetAUser(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "id")

	id, _ := primitive.ObjectIDFromHex(_id)

	w.Header().Set("Content-Type", "application/json")
	var user model.User

	database.User.UserColl.FindOne(r.Context(), bson.D{{Key: "_id", Value: id}}, options.FindOne().SetProjection(bson.D{{Key: "password", Value: 0}, {Key: "createdAt", Value: 0}})).Decode(&user)

	if user.ID == primitive.NilObjectID {
		log.Printf("Invalid get user request.")
		http.Error(w, "No user found with the given user id.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func AddAUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	user.Password, err = HashPassword(user.Password)

	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	newUser, err := database.User.UserColl.InsertOne(r.Context(), user)

	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	user.ID = newUser.InsertedID.(primitive.ObjectID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	// Producing user
	go producers.ProduceMessage(user, producers.CREATED)

}

func VerifyUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		log.Panic(err)
		return
	}

	password := user.Password
	email := user.Email

	err = database.User.UserColl.FindOne(r.Context(), bson.D{{Key: "email", Value: email}}).Decode(&user)

	if err != nil {

		http.Error(w, "Can't process your request at this moment.", http.StatusInternalServerError)
		log.Printf("Unsuccessful Login: %s\n", err)
		return
	}

	if user.ID == primitive.NilObjectID {

		http.Error(w, "No user found with the provided email id.", http.StatusBadRequest)
		log.Printf("No user found error %s", email)
		return
	}

	if ok := CheckPasswordHash(password, user.Password); !ok {

		http.Error(w, "E-mail or passwor is incorrect. Please try again.", http.StatusBadRequest)
		log.Printf("Unsuccessful Login %s\n", user.ID)
		return
	}

	// removing security concered details
	user.Password = ""
	user.CreatedAt = 0

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}
