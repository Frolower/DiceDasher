package rooms 

import (
	"encoding/json"
)

type Response struct {
	Room_id string `json:"room_id"`
	Player_id string `json:"player_id"`
}


func (this Response) JSON() []byte {
	res, _ := json.Marshal(this)
	// TODO: Handle error 
	return res
}