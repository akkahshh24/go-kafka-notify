package models

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	Store *NotificationStore
}

//Setup, Cleanup and ConsumeClaim are methods of the sarama.ConsumeGroupHandler interface

// Setup implements sarama.ConsumerGroupHandler.
func (*Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		userID := string(msg.Key)

		var notification Notification
		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}
		c.Store.Add(userID, notification)
		session.MarkMessage(msg, "")
	}
	return nil
}
