package models

import (
	"time"
)

type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	ReceiverID     int       `json:"receiver_id"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}

type MessageResponse struct {
	Message Message `json:"message"`
	Status  string  `json:"status"`
}
