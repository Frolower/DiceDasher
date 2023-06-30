package actions

import (
	// "dicedasher/rooms"
	"dicedasher/st"
	"encoding/json"
	"math/rand"
	"time"
)

func RollDice(size int) []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(size) + 10)
	for i := 0; i < numIter; i++ {
		result = append(result, timeRand.Intn(4)+1)
	}
	return result
}

func Roll(request st.Request) {
	r := RollDice(100)
	response := st.Response{
			Room_id: request.Room_id,
			Player_id: request.Player_id,
			Action: request.Action,
	}
	data, _ := json.Marshal(r)
	response.Data = string(data)
	send(response)
	// r, _:=json.Marshal(roll.RollDice(100))
	// players := RoomStorage[this.room_id].players // get all players in room
	// for i := 0; i < len(players); i++ {  // Send a message to each player
	// 	res := &Response{
	// 		Room_id: this.room_id,
	// 		Player_id: this.player_id,
	// 		Action: this.action,
	// 	}
	// 	res.Data= make(map[string]string)
	// 	res.Data["result"] = string(r)
		
	// }
}
