package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/RushinShah22/e-commerce-micro/services/products/pkg/database"
	model "github.com/RushinShah22/e-commerce-micro/services/products/pkg/models"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// partition number
const (
	CREATED = iota // partition 0
)

func SetupConsumer(groupID string, topics []string, topicPartition *[]kafka.TopicPartition, callback func(*kafka.Message)) {
	var err error
	Consumer, err := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": "kafka1:19092", "group.id": groupID, "auto.offset.reset": "smallest"})

	if err != nil {
		panic(err)
	}

	err = Consumer.SubscribeTopics(topics, nil)

	if err != nil {
		panic(err)
	}

	// will onyly consume from a assigned partition
	err = Consumer.Assign(*topicPartition)
	if err != nil {
		panic(err)
	}

	for {
		ev := Consumer.Poll(100) // checking for message every 100 ms
		switch e := ev.(type) {
		case *kafka.Message:
			// application-specific processing
			go callback(e)

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			os.Exit(1)
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}

}

// This function is called to consume a new order.
func OrderCallback(msg *kafka.Message) {
	var data interface{}
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		panic(err)
	}
	var productID primitive.ObjectID
	var quantityUsed int

	// moving fields around for struct field matching with the model
	if order, ok := data.(map[string]interface{}); ok {
		productID, _ = primitive.ObjectIDFromHex(order["productID"].(string))
		quantityUsed = int(order["quantity"].(float64))
	}
	switch msg.TopicPartition.Partition {
	case CREATED:
		var product model.Product
		database.Product.ProductColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: productID}}).Decode(&product)
		database.Product.ProductColl.FindOneAndUpdate(context.TODO(), bson.D{{Key: "_id", Value: productID}}, bson.M{"$set": bson.M{"quantity": product.Quantity - quantityUsed}})
		log.Printf("Consumed new product %s", productID)
	}

}

// This function is called when a new user is consumed
func UserCallback(msg *kafka.Message) {
	var userJson interface{}
	json.Unmarshal(msg.Value, &userJson)

	var seller model.Seller

	// moving field data around to make it compatible with the seller model.
	if data, ok := userJson.(map[string]interface{}); ok {

		if data["role"] != "seller" {
			return
		}
		data["userID"], _ = primitive.ObjectIDFromHex(data["id"].(string)) // since id -> userID
		data["id"] = ""

		tmpJson, _ := json.Marshal(data)
		json.Unmarshal(tmpJson, &seller)
	}

	insertedData, err := database.Product.SellerColl.InsertOne(context.Background(), seller)

	if err != nil {
		panic(err)
	}
	seller.ID = insertedData.InsertedID.(primitive.ObjectID)
	log.Printf("Consumed a new seller in products %s\n", seller.ID)
}
