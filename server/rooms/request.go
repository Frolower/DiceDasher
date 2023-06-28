package rooms 

import (
	"encoding/json"
)

type Request struct {
	// Structure that describes request from a client
	room_id string 
	player_id string 
	action string 
}

func (this *Request) fromJSON(JSON []byte) {
	reqMap := map[string]string{}
	json.Unmarshal(JSON, &reqMap)
	this.room_id = reqMap["room_id"]
	this.player_id = reqMap["player_id"]
	this.action = reqMap["action"]
}

func (this Request) handleAction() {
	if this.action == "ping" { // Check action type
		players := RoomStorage[this.room_id].players // get all players in room
		for i := 0; i < len(players); i++ {  // Send a message to each player
			clients[players[i]].WriteMessage(1, []byte("ping"))
		}
	}
}