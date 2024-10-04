package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/controllers"
	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/database"
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

	// ORDER router
	root.Route("/orders", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from Orders"))
		})
		r.Post("/", controllers.CreateOrder)
		r.Get("/", controllers.GetAllOrders)
	})

	// start server
	log.Printf("Server started on %s", port)
	if err := http.ListenAndServe(":"+port, root); err != nil {
		panic(err)
	}

}
