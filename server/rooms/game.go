package rooms

import (
	"dicedasher/actions"
	"dicedasher/st"
	"dicedasher/storage"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Game(c *gin.Context) {

	player_id := c.DefaultQuery("player_id", "nil")

	// TODO: Handle no params error
	
	ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil) // Upgrading HTTP connection to websocket
	storage.Clients[player_id] = ws
	for {
		fmt.Println(storage.Clients)
		fmt.Println(storage.RoomStorage)
		_, message, err := ws.ReadMessage()
		if err != nil { // If connection closed
			delete(storage.Clients, player_id) // Delete connection from clients
			break
		}
		req := st.Request{}
		req.FromJSON(message)
		if req.Action == "roll" {
			fmt.Println(req)
			actions.Roll(req)
		}
	}
}