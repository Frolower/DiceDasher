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

		ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil) // Upgrading HTTP connection to websocket
		// TODO: handle error

		clients[player_id] = ws // Add connection to map 

		for { // Reading incoming messages
			res := &Response{Room_id: room_id, Player_id: player_id} // Default info answer
			ws.WriteMessage(1, res.JSON())
	
			_, message, _ := ws.ReadMessage() // TODO: Handle error 
			
			req := Request{}
			req.fromJSON(message)
			req.handleAction()	
		}

	}
	c.JSON(http.StatusNoContent, gin.H{
		"message" : "no_room_id",
	})	
}