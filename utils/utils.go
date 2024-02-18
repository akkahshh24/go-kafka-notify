package utils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/akkahshh24/go-kafka-notify/pkg/models"
	"github.com/gin-gonic/gin"
)

var ErrUserNotFoundInProducer = errors.New("user not found")
var ErrNoMessagesFound = errors.New("no messages found")

func GetIDFromRequest(ctx *gin.Context, formValue string) (int, error) {
	id, err := strconv.Atoi(ctx.PostForm(formValue))
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID from form value %s: %w", formValue, err)
	}

	return id, nil
}

func FindUserByID(id int, users []models.User) (*models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, ErrUserNotFoundInProducer
}

func GetUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessagesFound
	}

	return userID, nil
}
