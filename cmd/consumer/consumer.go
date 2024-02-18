package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/akkahshh24/go-kafka-notify/pkg/models"
	"github.com/akkahshh24/go-kafka-notify/utils"
	"github.com/gin-gonic/gin"
)

const (
	KafkaServerAddress = "localhost:9092"
	ConsumerGroup      = "notifications-group"
	ConsumerTopic      = "notifications"
	ConsumerPort       = ":8082"
)

func handleNotifications(ctx *gin.Context, store *models.NotificationStore) {
	userID, err := utils.GetUserIDFromRequest(ctx)
	if err != nil {
		// 404 Not Found
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	notifications := store.Get(userID)
	if len(notifications) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message":       "No notifications found for user",
			"notifications": []models.Notification{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()

	// connect to the broker running on localhost:9092
	consumerGroup, err := sarama.NewConsumerGroup([]string{KafkaServerAddress}, ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	return consumerGroup, nil
}

func setupConsumerGroup(ctx context.Context, store *models.NotificationStore) {
	consumerGroup, err := initializeConsumerGroup()
	if err != nil {
		log.Printf("consumer group initialization error :%v", err)
	}
	defer consumerGroup.Close()

	consumer := models.Consumer{
		Store: store,
	}

	for {
		err := consumerGroup.Consume(ctx, []string{ConsumerTopic}, &consumer)
		if err != nil {
			log.Printf("error from consumer: %v", err)
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func main() {
	store := models.NotificationStore{
		Data: make(models.UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go setupConsumerGroup(ctx, &store)
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		handleNotifications(ctx, &store)
	})

	log.Printf("Kafka Consumer Group: %s started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	if err := router.Run(ConsumerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
