package rooms

import (
	roll "dicedasher/rooms/userActions"
	"encoding/json"
	"fmt"
)

type Request struct {
	// Structure that describes request from a client
	room_id string 
	player_id string 
	action string 
	data string
}

func (this *Request) fromJSON(JSON []byte) {
	reqMap := map[string]string{}
	json.Unmarshal(JSON, &reqMap)
	this.room_id = reqMap["room_id"]
	this.player_id = reqMap["player_id"]
	this.action = reqMap["action"]
	this.data = reqMap["data"]
}

func (this Request) handleAction() {
	if RoomStorage[this.room_id].isPlayerConnected(this.player_id) {
		fmt.Println(RoomStorage)
		fmt.Println(clients)
		dataJSON := map[string]string{}
		json.Unmarshal([]byte(this.data), &dataJSON)
		switch this.action {
		case "roll":
			r, _:=json.Marshal(roll.RollD100())
			players := RoomStorage[this.room_id].players // get all players in room
			for i := 0; i < len(players); i++ {  // Send a message to each player
				res := &Response{
					Room_id: this.room_id,
					Player_id: this.player_id,
					Action: this.action,
				}
				res.Data= make(map[string]string)
				res.Data["result"] = string(r)
				clients[players[i]].WriteMessage(1, res.JSON())
			}
			fmt.Println("message")
			

		default:
			
			fmt.Println("ping")
		}
	} else {
		fmt.Println("Player not connected")
		// TODO: Handle error (player not connected)
	}
}