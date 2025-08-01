package services

import (
	"SimpleChat/src/models"
	"time"
)

type MessageService struct {
	messages map[int]models.Message
	nextID   int
}

func NewMessageService() *MessageService {
	return &MessageService{
		messages: make(map[int]models.Message),
		nextID:   1,
	}
}

func (ms *MessageService) Create(message models.Message) (*models.Message, error) {
	message.ID = ms.nextID
	message.CreatedAt = time.Now()
	ms.messages[message.ID] = message
	ms.nextID = ms.nextID + 1

	return &message, nil

}
