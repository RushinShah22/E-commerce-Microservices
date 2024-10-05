package producers

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	CREATED = iota
)

type OrderProducer struct {
	Producer *kafka.Producer
	Topic    string
}

var Order OrderProducer

func SetupProducer() {
	var err error
	Order.Producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka1:19092"})
	Order.Topic = "orders"
	if err != nil {
		panic(err)
	}
}

func ProduceMessage(msg interface{}, partition int32) {
	msgBytes, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	err = Order.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &Order.Topic, Partition: partition},
		Value:          msgBytes,
	}, nil)

	if err != nil {
		panic(err)
	}
	log.Println("Produced a new order")

}
