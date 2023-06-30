package rooms

import (
	"dicedasher/st"
	"dicedasher/storage"

	"github.com/gin-gonic/gin"
)

func JoinRoom(c *gin.Context) {
	room_id := c.DefaultQuery("room_id", "nil")
	room, ok := storage.RoomStorage[room_id]
	
	if room_id == "nil" || !ok { // No room / no room_id error 
		c.JSON(400, st.Response{
			Room_id: room_id,
			Player_id: "",
			Action: "connect",
			Status: "error",
			Data: "{}",
		})	
		return 
	}

	player_id := generateID()
	room.Players = append(room.Players, player_id)
	storage.RoomStorage[room_id] = room 

	c.JSON(200, st.Response{
		Room_id: room_id,
		Player_id: player_id,
		Action: "connect",
		Status: "success",
		Data: "{}",
	})	
	return 
}