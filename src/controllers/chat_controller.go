package controllers

import (
	"SimpleChat/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	var body *models.SendMessage
	c.Bind(&body)

	c.JSON(http.StatusOK, gin.H{
		"data": &body,
	})
}

func GetMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "mensagem 1",
	})
}
