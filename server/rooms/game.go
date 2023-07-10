package rooms

import (
	"dicedasher/actions"
	"dicedasher/st"
	"dicedasher/storage"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Game(c *gin.Context) {

	player_id := c.DefaultQuery("player_id", "nil") // Get URL player_id param
	room_id := c.DefaultQuery("room_id", "nil")
	room, exists := storage.RoomStorage[room_id]

	if room_id == "nil" { // No player_id param
		// BAD REQUEST
		c.JSON(400, nil)
		return 
	}
	if !exists {
		// ACCESS DENIED
		c.JSON(401, nil)
		return 
	}
	if player_id == "nil" {
		player_id = generateID()
	}
	
	room.Players = append(room.Players, player_id)
	storage.RoomStorage[room_id] = room
	
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil) // Upgrading HTTP connection to websocket
	if err != nil { // Websocket error 
		c.JSON(500, nil)
	}
	
	storage.Clients[player_id] = ws
	actions.Send(st.Response{
		Room_id: room_id,
		Player_id: player_id,
		Action: "player_joined",
		Status: "success",
		Data: "{ players: [" + strings.Join(room.Players, ", ") + "]}",
	})
	for {
		mt, message, err := ws.ReadMessage()
	
		if err != nil || mt == websocket.CloseMessage { // error/closed conection
			_, ok := storage.Clients[player_id]
			if ok {
				delete(storage.Clients, player_id) 
			}
			break
		}

		req := st.Request{}
		req.Room_id = room_id
		req.Player_id = player_id
		result := req.FromJSON(message)
		if result == false { // bad JSON error 
			ws.WriteMessage(1, st.Response{
				Action: "error",
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