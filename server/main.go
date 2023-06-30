package main

import (
	rooms "dicedasher/rooms"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	server := gin.Default()
	server.Use(cors.Default())

	server.GET("/newroom", rooms.NewRoom)

	server.GET("/joinroom", rooms.JoinRoom)

	server.GET("/game", rooms.Game)

	server.Run()
}