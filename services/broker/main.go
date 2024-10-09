package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	// We create new topic using data from Args of the command
	if len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr,
			"Usage: %s <bootstrap-servers> <topic> <partition-count> <replication-factor>\n",
			os.Args[0])
		os.Exit(1)
	}

	bootstrapServers := os.Args[1]            // Get the server URL
	topic := os.Args[2]                       // Name of the topic we want to create
	numParts, err := strconv.Atoi(os.Args[3]) // Number of partition in this new topic.

	if err != nil {
		fmt.Printf("Invalid partition count: %s: %v\n", os.Args[3], err)
		panic(err)
	}
	replications, err := strconv.Atoi(os.Args[4])

	if err != nil {
		fmt.Printf("Invalid replication factor: %s: %v\n", os.Args[4], err)
		panic(err)
	}

	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	}
	adminClient, err := kafka.NewAdminClient(config)

	if err != nil {
		panic(err)
	}

	topicSpecs := kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     numParts,
		ReplicationFactor: replications,
	}

	results, err := adminClient.CreateTopics(context.Background(), []kafka.TopicSpecification{topicSpecs})

	if err != nil {
		fmt.Printf("Failed to create topic: %v\n", err)
		panic(err)
	}

	// Print results
	for _, result := range results {
		log.Printf("%s\n", result)
	}

	adminClient.Close()

}
