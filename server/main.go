package main

import (
	"net/http"
	rooms "dicedasher/rooms"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg" : "hello",
		})
	})

	server.GET("/newroom", rooms.NewRoom)

	server.GET("/joinroom", rooms.JoinRoom)

	server.Run()
}