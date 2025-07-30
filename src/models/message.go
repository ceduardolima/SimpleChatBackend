package models

import "time"

type Message struct {
	ID             int
	ConversationID int
	SenderID       int
	Message        string
	CreatedAt      time.Time
}
