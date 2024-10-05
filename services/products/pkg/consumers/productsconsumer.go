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

const (
	CREATED = iota
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
			type Order struct {
				ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
				ProductID   primitive.ObjectID `bson:"productID,omitempty" json:"productID,omitempty"`
				Quantity    int                `bson:"quantity,omitempty" json:"quantity,omitempty"`
				UserID      primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
				DateOfOrder primitive.DateTime `bson:"dateOfOrder,omitempty" json:"dateOfOrder,omitempty"`
				Status      string             `bson:"status,omitempty" json:"status,omitempty"`
			}
			var data Order
			err := json.Unmarshal(e.Value, &data)
			if err != nil {
				panic(err)
			}
			switch e.TopicPartition.Partition {
			case CREATED:
				var product model.Product
				database.Product.ProductColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: data.ProductID}}).Decode(&product)
				database.Product.ProductColl.FindOneAndUpdate(context.TODO(), bson.D{{Key: "_id", Value: data.ProductID}}, bson.M{"$set": bson.M{"quantity": product.Quantity - data.Quantity}})
				log.Printf("Consumed new product %s", data.ProductID)
			}

		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			os.Exit(1)
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}

}
