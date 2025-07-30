package router

import (
	"SimpleChat/src/services"
	"SimpleChat/src/websockets"

	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()
	cs := services.NewConversationService()
	ms := services.NewMessageService()
	ws := websockets.NewWSServer(cs, ms)

	go ws.HandleMessages()

	router.GET("/chat", ws.HandleConnections)

	return router
}
