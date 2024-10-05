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

func SetupConsumer(groupID string, topics []string, topicPartition *[]kafka.TopicPartition) {
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
			// application-specific processing
			var data model.Catalog
			err := json.Unmarshal(e.Value, &data)
			if err != nil {
				panic(err)
			}
			switch e.TopicPartition.Partition {
			case CREATED:
				insertedData, err := database.Order.CatalogColl.InsertOne(context.TODO(), data)
				if err != nil {
					panic(err)
				}
				log.Printf("Consumed new product %s", insertedData.InsertedID)
			case UPDATED:
				database.Order.CatalogColl.FindOneAndReplace(context.TODO(), bson.D{{Key: "productID", Value: data.ProductID}}, data, options.FindOneAndReplace().SetReturnDocument(options.After))
				log.Printf("Consumed Updated product %s", data.ProductID)
			}

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			os.Exit(1)
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}

}
