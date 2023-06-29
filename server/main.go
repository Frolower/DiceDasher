package main

import (
	rooms "dicedasher/rooms"
	"net/http"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	server := gin.Default()
	server.Use(cors.Default())
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg" : "hello",
		})
	})

	server.GET("/newroom", rooms.NewRoom)

	server.GET("/joinroom", rooms.JoinRoom)

	server.GET("/game", rooms.Game)

	server.Run()
}