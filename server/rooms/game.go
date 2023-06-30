package rooms

import (
	"dicedasher/actions"
	"dicedasher/st"
	"dicedasher/storage"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Game(c *gin.Context) {

	player_id := c.DefaultQuery("player_id", "nil") // Get URL player_id param
	if player_id == "nil" { // No player_id param
		c.JSON(400, nil)
	}
	
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil) // Upgrading HTTP connection to websocket
	if err != nil { // Websocket error 
		c.JSON(500, nil)
	}

	storage.Clients[player_id] = ws
	for {
		mt, message, err := ws.ReadMessage()

		if err != nil || mt == websocket.CloseMessage { // error/closed conection
			delete(storage.Clients, player_id) 
			break
		}

		req := st.Request{}
		result := req.FromJSON(message)
		if result == false || player_id != req.Player_id{ // bad JSON error / different ids 
			ws.WriteMessage(1, st.Response{
				Player_id: player_id,
				Action: "undefined",
				Status: "bad_request",
			}.JSON())
			return 
		}
		

		switch req.Action { // Action handlers
			case "roll":
				go actions.Roll(req)
			default:
				fmt.Println(req)
		}
	}
}