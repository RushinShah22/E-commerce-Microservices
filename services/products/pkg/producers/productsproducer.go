package producers

import (
	"encoding/json"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const (
	CREATED = iota
	UPDATED
)

type ProductProducer struct {
	Producer *kafka.Producer
	Topic    string
}

var Product ProductProducer

func SetupProducer() {
	var err error
	Product.Producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092",
		"acks": "all"})
	Product.Topic = "products"
	if err != nil {
		panic(err)
	}
}

func ProduceMessage(msg interface{}, partition int32) {
	msgBytes, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	err = Product.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &Product.Topic, Partition: partition},
		Value:          msgBytes,
	}, nil)

	if err != nil {
		panic(err)
	}

}
