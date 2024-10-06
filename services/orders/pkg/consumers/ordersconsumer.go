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

const (
	CREATED = iota
	UPDATED
)

var Order *kafka.Consumer

func SetupConsumer(groupID string, topics []string, topicPartition *[]kafka.TopicPartition, callback func(*kafka.Message)) {
	var err error
	Order, err = kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": "kafka1:19092", "group.id": groupID, "auto.offset.reset": "smallest"})

	if err != nil {
		panic(err)
	}

	err = Order.SubscribeTopics(topics, nil)

	if err != nil {
		panic(err)
	}
	err = Order.Assign(*topicPartition)
	if err != nil {
		panic(err)
	}

	for {
		ev := Order.Poll(100)
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

func ProductsCallback(msg *kafka.Message) {
	var data model.Catalog
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		panic(err)
	}
	switch msg.TopicPartition.Partition {
	case CREATED:
		insertedData, err := database.Order.CatalogColl.InsertOne(context.Background(), data)
		if err != nil {
			panic(err)
		}
		log.Printf("Consumed new product %s", insertedData.InsertedID)
	case UPDATED:
		database.Order.CatalogColl.FindOneAndReplace(context.TODO(), bson.D{{Key: "productID", Value: data.ProductID}}, data, options.FindOneAndReplace().SetReturnDocument(options.After))
		log.Printf("Consumed Updated product %s", data.ProductID)
	}

}

func UserCallback(msg *kafka.Message) {
	fmt.Println("yess")
	var userJson interface{}
	err := json.Unmarshal(msg.Value, &userJson)

	if err != nil {
		panic(err)
	}
	var user model.User

	if data, ok := userJson.(map[string]interface{}); ok {
		data["userID"] = data["_id"]
		data["_id"] = ""
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
