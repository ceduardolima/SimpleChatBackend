package router

import (
	"SimpleChat/src/controllers"

	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/chat", controllers.SendMessage)
	router.GET("/chat", controllers.GetMessage)

	return router
}
