package models

type SendMessage struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
}

type GetMessage struct {
	UserID int `json:"user_id"`
}
