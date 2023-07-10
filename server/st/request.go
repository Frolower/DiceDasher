package st

import (
	"github.com/tidwall/gjson"
)

type Request struct {
	// Structure that describes request from a client
	Action    string 
	Room_id   string 
	Player_id string 
	Data gjson.Result
}

func (g *Request) FromJSON(JSON []byte) bool{
	request := gjson.Parse(string(JSON))
	g.Action = request.Get("action").String()
	g.Data = request.Get("data")
	if g.Action == "" {
		return false 
	}
	return true 
}
