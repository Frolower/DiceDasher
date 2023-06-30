package st

import (
	"encoding/json"
)

type Request struct {
	// Structure that describes request from a client
	Room_id string 
	Player_id string 
	Action string 
	Data string
}

func (g *Request) FromJSON(JSON []byte) {
	var requestMap map[string]string
	json.Unmarshal(JSON, &requestMap)
	g.Room_id = requestMap["room_id"]
	g.Player_id = requestMap["player_id"]
	g.Action = requestMap["action"]
	g.Data = requestMap["data"]
}

// func (this Request) handleAction() {
// 	if storage.RoomStorage[this.room_id].isPlayerConnected(this.player_id) {
// 		dataJSON := map[string]string{}
// 		json.Unmarshal([]byte(this.data), &dataJSON)
// 		switch this.action {
// 		case "roll":
// 			actions.Roll()
// 		default:
// 			fmt.Println("ping")
// 		}
// 	} else {
// 		fmt.Println("Player not connected")
// 		// TODO: Handle error (player not connected)
// 	}
// }