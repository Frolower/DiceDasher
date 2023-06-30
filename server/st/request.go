package st

import (
	"github.com/tidwall/gjson"
)

type Request struct {
	// Structure that describes request from a client
	Room_id string 
	Player_id string 
	Action string 
	Data gjson.Result
}

func (g *Request) FromJSON(JSON []byte) bool{
	request := gjson.Parse(string(JSON))
	g.Room_id = request.Get("room_id").String()
	g.Player_id = request.Get("player_id").String()
	g.Action = request.Get("action").String()
	g.Data = request.Get("data")
	if g.Room_id == "" || g.Player_id == "" || g.Action == "" {
		return false 
	}
	return true 
	// var requiredFields =  []string{"room_id", "player_id", "action", "data"} // Required field for every Request
	// var requestMap map[string]string 
	// err := json.Unmarshal(JSON, &requestMap) 
	// requestMapKeys := maps.Keys(requestMap)
	// fmt.Println(requestMapKeys)

	// fmt.Println(gjson.Get(string(JSON), "data"))
	// if err != nil || len(requiredFields) != len(requestMapKeys) { 
	// 	return false 
	// }

	// for i:=0; i < len(requiredFields); i++ {
	// 	if !slices.Contains(requestMapKeys, requiredFields[i]) {
	// 		return false
	// 	}
	// }
	
	// g.Room_id = requestMap["room_id"]
	// g.Player_id = requestMap["player_id"]
	// g.Action = requestMap["action"]
	// // g.Data = requestMap["data"]
	return true 
}
