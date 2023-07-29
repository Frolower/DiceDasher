package rooms

import (
	"dicedasher/actions"
	"dicedasher/st"
	"dicedasher/storage"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func playerList(room st.Room) string {
	data, _ := json.Marshal(map[string]string{
		"players": "[" + strings.Join(room.Players, ", ") + "]",
	})
	return string(data)
}

func Game(c *gin.Context) {

	player_id := c.DefaultQuery("player_id", "nil") // Get URL player_id param
	room_id := c.DefaultQuery("room_id", "nil")
	room, exists := storage.RoomStorage[room_id]

	if room_id == "nil" { // No room_id param
		// BAD REQUEST
		c.JSON(400, nil)
		return
	}
	if !exists {
		// ACCESS DENIED
		c.JSON(401, nil)
		return 
	}
	if player_id == "nil" || player_id == ""{
		player_id = generateID()
		room.Players = append(room.Players, player_id)
	} else {
		if player_id != room.Master {
			c.JSON(401, nil)
			return
		}
	}

	storage.RoomStorage[room_id] = room

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil) // Upgrading HTTP connection to websocket
	if err != nil {                                       // Websocket error
		c.JSON(500, nil)
	}

	storage.Clients[player_id] = ws

	connected := false

	for {
		fmt.Println(storage.RoomStorage[room_id])
		if !storage.RoomStorage[room_id].IsOpened && !connected {
			ws.WriteMessage(1, st.Response{
				Room_id:   room_id,
				Player_id: player_id,
				Action:    "connect",
				Status:    "room_closed",
			}.JSON())
			actions.TryCloseNormally(ws) // forces to close connection
			return
		}
		if !connected {

			actions.Send(st.Response{ // Send message when player is connected
				Room_id:   room_id,
				Player_id: player_id,
				Action:    "connect",
				Status:    "success",
				Data:      playerList(room),
			})
			connected = true
		}
		mt, message, err := ws.ReadMessage()

		if err != nil || mt == websocket.CloseMessage { // error/closed conection
			_, ok := storage.Clients[player_id]
			if ok {
				delete(storage.Clients, player_id)
			}
			room.Players = room.RemovePlayer(player_id)
			fmt.Println(room)
			actions.Send(st.Response{ // Send message when player is connected
				Room_id:   room_id,
				Player_id: player_id,
				Action:    "disconnect",
				Status:    "success",
				Data:      playerList(room),
			})
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

		correct := true
		switch req.Action { // Action handlers
		case "open_room":
			correct = actions.Lock(req)
		case "close_room":
			correct = actions.Lock(req)
		case "roll":
			correct = actions.Roll(req)
		default:
			fmt.Println(req)
		}

		if !correct {
			ws.WriteMessage(1, st.Response{
				Player_id: player_id,
				Room_id:   room_id,
				Action:    req.Action,
				Status:    "bad_request",
			}.JSON())
		}
	}
}
