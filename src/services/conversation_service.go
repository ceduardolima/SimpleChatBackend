package services

import (
	"SimpleChat/src/models"
	"errors"
	"log"
)

type ConversationService struct {
	conversations []models.Conversation
}

func NewConversationService() *ConversationService {
	return &ConversationService{conversations: []models.Conversation{}}
}

func (cs *ConversationService) Create(conversation models.Conversation) (*models.Conversation, error) {

	for _, c := range cs.conversations { // Usando for...range é mais idiomático
		if c.SenderID == conversation.SenderID && c.ReceiverID == conversation.ReceiverID {
			log.Println("Conversation conflict")
			return nil, errors.New("Conversation Conflict")
		}
	}

	conversation.ID = 1
	cs.conversations = append(cs.conversations, conversation)

	return &conversation, nil
}
