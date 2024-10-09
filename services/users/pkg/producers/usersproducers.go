package producers

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Parition number
const (
	CREATED = iota // a new user is created on parition 0
)

type UserProducer struct {
	Producer *kafka.Producer
	Topic    string
}

var User UserProducer

func SetupProducer() {
	var err error
	User.Producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka1:19092"})
	User.Topic = "users"
	if err != nil {
		panic(err)
	}
}

func ProduceMessage(msg interface{}, partition int32) {
	msgBytes, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	err = User.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &User.Topic, Partition: partition},
		Value:          msgBytes,
	}, nil)

	if err != nil {
		panic(err)
	}
	log.Println("Produced a User.")
}
