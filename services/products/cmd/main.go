package main

import (
	"context"
	"log"
	"net/http"
	"os"

	consumer "github.com/RushinShah22/e-commerce-micro/services/products/pkg/consumers"
	"github.com/RushinShah22/e-commerce-micro/services/products/pkg/controllers"
	"github.com/RushinShah22/e-commerce-micro/services/products/pkg/database"
	"github.com/RushinShah22/e-commerce-micro/services/products/pkg/producers"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	port := os.Getenv("PORT")

	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	// DB setup
	database.ConnToDB(uri)
	defer database.Product.Client.Disconnect(context.TODO())

	// Create the Router entry point
	root := chi.NewRouter()

	// PRODUCT router
	root.Route("/products", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello From Products."))
		})
		r.Get("/", controllers.GetAllProducts)
		r.Post("/", controllers.CreateProduct)
		r.Get("/{id}", controllers.GetAProduct)
		r.Patch("/{id}", controllers.UpdateProduct)

	})

	// Setup consumers
	// Consumer for orders produced.
	orderTopic := "orders"
	go consumer.SetupConsumer("orders", []string{"orders"}, &[]kafka.TopicPartition{{Topic: &orderTopic, Partition: consumer.CREATED}}, consumer.OrderCallback)

	// Consumer for new Seller created
	sellerTopics := "users"
	go consumer.SetupConsumer("sellers", []string{"users"}, &[]kafka.TopicPartition{{Topic: &sellerTopics, Partition: consumer.CREATED}}, consumer.UserCallback)

	// Create Producer
	producers.SetupProducer()
	defer producers.Product.Producer.Close()

	// start server
	log.Printf("Server started on %s", port)
	if err := http.ListenAndServe(":"+port, root); err != nil {
		panic(err)
	}

}
