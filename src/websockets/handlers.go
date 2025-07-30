package websockets

import (
	"SimpleChat/src/models"
	"SimpleChat/src/services"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	clients             map[*websocket.Conn]bool
	broadcast           chan []byte
	register            chan *websocket.Conn
	unregister          chan *websocket.Conn
	conversationService *services.ConversationService
	messageService      *services.MessageService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWSServer(conversationService *services.ConversationService, messageService *services.MessageService) *WSServer {
	return &WSServer{
		clients:             make(map[*websocket.Conn]bool),
		broadcast:           make(chan []byte),
		register:            make(chan *websocket.Conn),
		unregister:          make(chan *websocket.Conn),
		conversationService: conversationService,
		messageService:      messageService,
	}
}

func (ws *WSServer) HandleMessages() {
	for {
		select {
		case conn := <-ws.register:
			ws.clients[conn] = true
		case conn := <-ws.unregister:
			if _, ok := ws.clients[conn]; ok {
				delete(ws.clients, conn)
				conn.Close()
			}
		case message := <-ws.broadcast:
			for range ws.clients {
				log.Println(message)
			}
		}
	}
}

func (ws *WSServer) HandleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to ws: %v\n", err)
		return
	}

	defer conn.Close()

	log.Printf("Client connected: %s", conn.RemoteAddr())

	ws.register <- conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			ws.unregister <- conn
			break
		}

		ws.broadcast <- msg

		var request map[string]any
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Printf("Error parsins message: %v\n", err)
			continue
		}

		log.Printf("message: %v\n", request)

		action, ok := request["action"].(string)
		if !ok {
			log.Println("Invalid action")
			continue
		}

		switch action {
		case "create_conversation":
			err := ws.createConversation(request)
			if err != nil {
				log.Printf("Create conversation error: %v\n", err)
			}
		case "send_message":
			response, err := ws.sendMessage(request)
			if err != nil {
				log.Fatalf("Send message error: %v\n", err)
				continue
			}
			err = conn.WriteMessage(websocket.TextMessage, response)
			if err != nil {
				log.Fatalf("Write message error: %v\n", err)
				continue
			}
			log.Println("Send message successfully")

		default:
			log.Println("Unknown action")
		}

	}
}

func (ws *WSServer) sendMessage(request map[string]any) ([]byte, error) {
	conversationIDF, ok1 := request["conversation_id"].(float64)
	senderIDF, ok2 := request["sender_id"].(float64)

	if !ok1 || !ok2 {
		return nil, errors.New("Invalid conversation or sender id")
	}

	conversationID := int(conversationIDF)
	senderID := int(senderIDF)

	message := models.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Message:        request["message"].(string),
	}

	createdMessage, err := ws.messageService.Create(message)
	if err != nil {
		return nil, fmt.Errorf("Error creating message: %v\n", err)
	}

	response := map[string]any{
		"status":  "success",
		"message": createdMessage,
	}

	responseMessage, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling response: %v\n", err)
	}

	return responseMessage, nil

}

func (ws *WSServer) createConversation(request map[string]any) error {

	senderIDF, ok1 := request["sender_id"].(float64)
	receiverIDF, ok2 := request["receiver_id"].(float64)

	if !ok1 || !ok2 {
		return errors.New("Invalid sender or receiver ID")
	}

	senderID := int(senderIDF)
	receiverID := int(receiverIDF)

	conversation := models.Conversation{
		SenderID:   senderID,
		ReceiverID: receiverID,
	}

	_, err := ws.conversationService.Create(conversation)
	if err != nil {
		return err
	}

	log.Printf("Conversation created!")
	return nil
}
