package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	database "github.com/RushinShah22/e-commerce-micro/services/users/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/users/pkg/models"
	"github.com/RushinShah22/e-commerce-micro/services/users/pkg/producers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	Message  string      `json:"message,omitempty"`
	Status   string      `json:"status,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Err      error       `json:"err,omitempty"`
	Verified bool        `json:"verified,omitempty"`
}

type UserResponseInterface interface {
	UserResponse
	String() string
}

func (res *UserResponse) String() string {
	jsonResponse, _ := json.Marshal(res)
	return string(jsonResponse)

}

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
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		panic(err)
	}

	user.Password, err = HashPassword(user.Password)

	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		panic(err)
	}
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

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

func VerifyUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)

	var response UserResponse

	if err != nil {
		response.Err = err
		response.Message = "Something went wrong."
		response.Status = "fail"

		http.Error(w, response.String(), http.StatusInternalServerError)
		panic(err)
	}

	password := user.Password
	email := user.Email

	err = database.User.UserColl.FindOne(r.Context(), bson.D{{Key: "email", Value: email}}).Decode(&user)

	if err != nil {
		response.Err = err
		response.Message = "Can't process your request at this moment."
		response.Status = "fail"

		http.Error(w, response.String(), http.StatusInternalServerError)
		log.Printf("Unsuccessful Login: %s\n", err)
		return
	}

	if user.ID == primitive.NilObjectID {
		response.Err = err
		response.Message = "No user found with the provided email id."
		response.Status = "fail"

		http.Error(w, response.String(), http.StatusBadRequest)
		log.Printf("No user found error %s", email)
		return
	}

	if ok := CheckPasswordHash(password, user.Password); !ok {
		response.Err = errors.New("Authentication Failed.")
		response.Message = "Email or Password is wrong."
		response.Status = "fail"

		http.Error(w, response.String(), http.StatusBadRequest)
		log.Printf("Unsuccessful Login %s\n", user.ID)
		return
	}
	w.WriteHeader(http.StatusOK)
	response.Message = "Login successful."
	response.Verified = true
	response.Status = "success"

	json.NewEncoder(w).Encode(response)
}
