package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RushinShah22/e-commerce-micro/services/users/pkg/controllers"
	client "github.com/RushinShah22/e-commerce-micro/services/users/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Connect to DB
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	client.ConnToDB(uri)
	// Create the entry point
	root := chi.NewRouter()

	// USER router
	root.Route("/user", func(r chi.Router) {
		r.Get("/", controllers.GetAllUser)
		r.Get("/{id}", controllers.GetAUser)
	})

	// start server
	log.Println("Server started on 3000")
	if err := http.ListenAndServe(":3000", root); err != nil {
		panic(err)
	}

}
