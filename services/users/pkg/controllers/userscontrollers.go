package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	client "github.com/RushinShah22/e-commerce-micro/services/users/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/users/pkg/models"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	coll := client.User.DB.Collection("users")

	cursor, err := coll.Find(r.Context(), bson.M{})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	var users []model.User

	if err := cursor.All(r.Context(), &users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	usersList, err := json.Marshal(users)
	w.Write(usersList)

}

func GetAUser(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "id")

	coll := client.User.DB.Collection("User")

	user := coll.FindOne(r.Context(), bson.M{"_id": _id})

	userJson, err := json.Marshal(user)

	if err != nil {
		panic("error")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}
