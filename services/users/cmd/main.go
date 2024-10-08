package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RushinShah22/e-commerce-micro/services/users/pkg/controllers"
	database "github.com/RushinShah22/e-commerce-micro/services/users/pkg/database"
	"github.com/RushinShah22/e-commerce-micro/services/users/pkg/producers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Connect to DB
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	port := os.Getenv("PORT")

	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	database.ConnToDB(uri)
	// Create the entry point
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
