package main

import (
	auth "dicedasher/auth"
	rooms "dicedasher/rooms"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	server := gin.Default()
	server.Use(cors.Default())

	server.GET("/newroom", rooms.NewRoom)

	server.GET("/game", rooms.Game)

	server.POST("/reg", auth.Reg)

	server.POST("/login", auth.Login)

	server.Run()
}