package rooms

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JoinRoom(c *gin.Context) {
	room_id := c.DefaultQuery("room_id", "nil")
	if room_id != "nil" {
		room := RoomStorage[room_id]

		// TODO:
		// Handle no room error 
		
		player_id := generateID()
		room.players = append(room.players, player_id)
		RoomStorage[room_id] = room 
		c.JSON(200, gin.H{
			"status" : "connected",
			"room_id" : room_id,
			"player_id" : player_id,
		})
	}
	c.JSON(http.StatusNoContent, gin.H{
		"status" : "no_room_id",
	})	
}