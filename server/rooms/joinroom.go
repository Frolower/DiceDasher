package rooms 

import (
	// "fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func JoinRoom(c *gin.Context) {
	room_id := c.DefaultQuery("room_id", "nil")
	if room_id != "nil" {
		room := RoomStorage[room_id]
		player_id := generateID()
		room.players = append(room.players, player_id)
		RoomStorage[room_id] = room 
		c.JSON(http.StatusOK, gin.H{
			"player_id" : player_id,
		})	
	}
	c.JSON(http.StatusNoContent, gin.H{
		"message" : "no_room_id",
	})
	
}