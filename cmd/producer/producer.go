package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/akkahshh24/go-kafka-notify/pkg/models"
	"github.com/akkahshh24/go-kafka-notify/utils"
	"github.com/gin-gonic/gin"
)

const (
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
	ProducerPort       = ":8080"
)

func sendKafkaMessage(ctx *gin.Context, producer sarama.SyncProducer, users []models.User, fromID, toID int) error {
	message := ctx.PostForm("message")

	fromUser, err := utils.FindUserByID(fromID, users)
	if err != nil {
		return err
	}

	toUser, err := utils.FindUserByID(toID, users)
	if err != nil {
		return err
	}

	notification := models.Notification{
		From:    fromUser,
		To:      toUser,
		Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(strconv.Itoa(toUser.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	partition, offset, err := producer.SendMessage(&msg)
	log.Printf("partition: %v, offset: %v", partition, offset)
	return err
}

func sendMessageHandler(producer sarama.SyncProducer, users []models.User) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fromID, err := utils.GetIDFromRequest(ctx, "fromID")
		if err != nil {
			// 400 Bad Request
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		toID, err := utils.GetIDFromRequest(ctx, "toID")
		if err != nil {
			// 400 Bad Request
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		err = sendKafkaMessage(ctx, producer, users, fromID, toID)
		if errors.Is(err, utils.ErrUserNotFoundInProducer) {
			// 404 Not Found
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		if err != nil {
			// 500 Internal Server Error
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully!"})
	}
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	// Producer receives an ack once the msg is stored in topic
	config.Producer.Return.Successes = true
	// connect to the Kafka broker running at localhost:9092
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

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/send", sendMessageHandler(producer, users))

	log.Printf("Kafka Producer started at http://localhost%s\n", ProducerPort)

	if err := router.Run(ProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
