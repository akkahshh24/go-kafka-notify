package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/akkahshh24/go-kafka-notify/pkg/models"
)

const (
	KafkaServerAddress = "localhost:9092"
)

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %v", err)
	}
	return producer, nil
}

func main() {
	users := []models.User{
		{ID: 1, Name: "Shyam"},
		{ID: 2, Name: "Raju"},
		{ID: 3, Name: "Anuradha"},
	}

	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer %v", err)
	}
	defer producer.Close()
}
