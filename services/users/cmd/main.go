package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RushinShah22/e-commerce-micro/services/users/pkg/controllers"
	database "github.com/RushinShah22/e-commerce-micro/services/users/pkg/database"
	"github.com/RushinShah22/e-commerce-micro/services/users/pkg/producers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to DB

	uri := os.Getenv("MONGODB_URI")
	port := os.Getenv("PORT")

	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	database.ConnToDB(uri)
	// Create the entry point

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}}, // Create index on the "email" field
		Options: options.Index().SetUnique(true),
	}
	ctx, canc := context.WithTimeout(context.Background(), 10*time.Second)
	defer canc()

	_, err := database.User.UserColl.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatal("Failed to create index:", err)
	}
	root := chi.NewRouter()

	// USER router
	root.Route("/users", func(r chi.Router) {
		r.Get("/", controllers.GetAllUser)
		r.Get("/{id}", controllers.GetAUser)
		r.Post("/", controllers.AddAUser)
		r.Post("/verify", controllers.VerifyUser)
	})

	// Setup up producers

	go producers.SetupProducer()

	// start server
	log.Printf("Server started on %s", port)
	if err := http.ListenAndServe(":"+port, root); err != nil {
		panic(err)
	}

}
