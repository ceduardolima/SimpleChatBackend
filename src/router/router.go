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
	us := services.NewUserService()

	ws := websockets.NewWSServer(cs, ms, us)

	go ws.HandleMessages()

	router.GET("/chat", ws.HandleConnections)

	return router
}
