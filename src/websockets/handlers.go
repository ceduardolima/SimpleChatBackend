package websockets

import (
	"SimpleChat/src/models"
	"SimpleChat/src/services"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type UserConn struct {
	conn *websocket.Conn
	user *models.User
}

type WSServer struct {
	clients             map[int]*UserConn
	broadcast           chan []byte
	register            chan *UserConn
	unregister          chan *UserConn
	conversationService *services.ConversationService
	messageService      *services.MessageService
	userService         *services.UserServices
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWSServer(conversationService *services.ConversationService, messageService *services.MessageService, userService *services.UserServices) *WSServer {
	return &WSServer{
		clients:             make(map[int]*UserConn),
		broadcast:           make(chan []byte),
		register:            make(chan *UserConn),
		unregister:          make(chan *UserConn),
		conversationService: conversationService,
		messageService:      messageService,
		userService:         userService,
	}
}

func (ws *WSServer) HandleMessages() {
	for {
		select {
		case userConn := <-ws.register:
			ws.clients[userConn.user.ID] = userConn

		case userConn := <-ws.unregister:
			if _, ok := ws.clients[userConn.user.ID]; ok {
				delete(ws.clients, userConn.user.ID)
				userConn.conn.Close()
			}

		case message := <-ws.broadcast:
			var request models.MessageResponse
			if err := json.Unmarshal(message, &request); err != nil {
				log.Printf("Error writing message: %v", err)
				continue
			}
			userConn := ws.clients[int(request.Message.ReceiverID)]
			err := userConn.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				userConn.conn.Close()
				delete(ws.clients, userConn.user.ID)

			}
		}
	}
}

func (ws *WSServer) HandleConnections(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")

	if authorization == "" {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	userID, err := strconv.Atoi(authorization[len("Bearer "):])

	if err != nil {
		log.Panicln(err)
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	user, err := ws.userService.GetUserById(userID)
	if err != nil {
		log.Panicln(err)
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to ws: %v\n", err)
		return
	}

	defer conn.Close()

	userConn := &UserConn{conn: conn, user: user}
	log.Printf("Client %s connected with id: %d", conn.RemoteAddr(), user.ID)

	ws.register <- userConn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			ws.unregister <- userConn
			break
		}

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

		ws.handleActions(action, request)
	}
}

func (ws *WSServer) handleActions(action string, request map[string]any) {
	switch action {

	case "create_conversation":
		err := ws.createConversation(request)
		if err != nil {
			log.Printf("Create conversation error: %v\n", err)
			return
		}

	case "send_message":
		err := ws.sendMessage(request)
		if err != nil {
			log.Fatalf("Send message error: %v\n", err)
			return
		}
		log.Println("Send message successfully")

	default:
		log.Println("Unknown action")
	}
}

func (ws *WSServer) sendMessage(request map[string]any) error {
	conversationIDF, ok1 := request["conversation_id"].(float64)
	senderIDF, ok2 := request["sender_id"].(float64)
	receiverIDF, ok3 := request["receiver_id"].(float64)

	if !ok1 || !ok2 || !ok3 {
		return errors.New("Invalid conversation or sender id")
	}

	conversationID := int(conversationIDF)
	senderID := int(senderIDF)
	receiverID := int(receiverIDF)

	message := models.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Message:        request["message"].(string),
	}

	createdMessage, err := ws.messageService.Create(message)
	if err != nil {
		return fmt.Errorf("Error creating message: %v\n", err)
	}

	response := map[string]any{
		"status":  "success",
		"message": createdMessage,
	}

	responseMessage, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("Error marshalling response: %v\n", err)
	}

	ws.broadcast <- responseMessage

	return nil

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
		UserID1: senderID,
		UserID2: receiverID,
	}

	_, err := ws.conversationService.Create(conversation)
	if err != nil {
		return err
	}

	log.Printf("Conversation created!")
	return nil
}
