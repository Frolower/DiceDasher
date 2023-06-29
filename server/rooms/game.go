package rooms

import (
	"github.com/gin-gonic/gin"
)

func Game(c *gin.Context) {

	player_id := c.DefaultQuery("player_id", "nil")

	// TODO: Handle no params error
	
	ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil) // Upgrading HTTP connection to websocket
	clients[player_id] = ws
	for {
		_, message, err := ws.ReadMessage()
		if err != nil { // If connection closed
			delete(clients, player_id) // Delete connection from clients
			break
		}
		req := Request{}
		req.fromJSON(message)
		req.handleAction()
	}
}