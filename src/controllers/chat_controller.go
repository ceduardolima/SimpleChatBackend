package controllers

import (
	"SimpleChat/src/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

func SendMessage(c *gin.Context) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	for {
		body := &models.SendMessage{}
		err := ws.ReadJSON(body)

		if err != nil {
			log.Panicln("read error:", err)

			response := models.MessageResponse{
				Status:  models.Error,
				Message: "Fail to deliver the message",
			}

			if err := ws.WriteJSON(response); err != nil {
				fmt.Println("write error:", err)
				break
			}
			continue
		}

		log.Printf("Received message from %d\n", body.SenderID)
		response := models.MessageResponse{
			Status: "ok",
		}
		if err := ws.WriteJSON(response); err != nil {
			fmt.Println("write error:", err)
			break
		}
	}
}

func GetMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "mensagem 1",
	})
}
