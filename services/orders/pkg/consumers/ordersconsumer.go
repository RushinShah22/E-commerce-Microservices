package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/database"
	"github.com/RushinShah22/e-commerce-micro/services/orders/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Partition number
const (
	CREATED = iota // parition 0 for create event
	UPDATED        // partition 1 for update event
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
	err = Consumer.Assign(*topicPartition)
	if err != nil {
		panic(err)
	}

	for {
		ev := Consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			go callback(e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			os.Exit(1)
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}

}

// This function is called when a product is consumed.
func ProductsCallback(msg *kafka.Message) {
	var productJson interface{}

	err := json.Unmarshal(msg.Value, &productJson)
	if err != nil {
		panic(err)
	}

	var product model.Catalog

	// changing data fields to match model.Catalog
	if data, ok := productJson.(map[string]interface{}); ok {
		data["productID"] = data["id"]
		data["id"] = ""
		dataJson, err := json.Marshal(data)

		if err != nil {
			panic(err)
		}
		json.Unmarshal(dataJson, &product)
	} else {
		panic(ok)
	}

	switch msg.TopicPartition.Partition {
	case CREATED: // A product is consumer
		insertedData, err := database.Order.CatalogColl.InsertOne(context.Background(), product)
		if err != nil {
			panic(err)
		}
		log.Printf("Consumed new product %s", insertedData.InsertedID)
	case UPDATED: // A product is updated
		database.Order.CatalogColl.FindOneAndReplace(context.TODO(), bson.D{{Key: "productID", Value: product.ProductID}}, product, options.FindOneAndReplace().SetReturnDocument(options.After))
		log.Printf("Consumed Updated product %s", product.ProductID)
	}

}

// This function is called when a new user event is consumed
func UserCallback(msg *kafka.Message) {
	var userJson interface{}
	err := json.Unmarshal(msg.Value, &userJson)

	if err != nil {
		panic(err)
	}
	var user model.User

	// changing data fields to match the model.User
	if data, ok := userJson.(map[string]interface{}); ok {
		data["userID"] = data["id"]
		data["id"] = ""
		tmpJson, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(tmpJson, &user)
	} else {
		panic(ok)
	}

	insertedUser, err := database.Order.UserColl.InsertOne(context.Background(), user)

	if err != nil {
		panic(err)
	}
	log.Printf("Consumed a new user. %s\n", insertedUser.InsertedID)

}
