package main

import (
	"context"
	"log"
	"net/http"
	"os"

	consumer "github.com/RushinShah22/e-commerce-micro/services/orders/pkg/consumers"
	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/controllers"
	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/database"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
	defer database.Order.Client.Disconnect(context.TODO())
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

	topic := "products"
	go consumer.SetupConsumer("products", []string{"products"}, &[]kafka.TopicPartition{
		{Topic: &topic, Partition: consumer.CREATED},
	})
	// start server
	log.Printf("Server started on %s", port)
	if err := http.ListenAndServe(":"+port, root); err != nil {
		panic(err)
	}

}
