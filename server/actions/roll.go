package actions

import (
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

func Roll() {
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
	// 	clients[players[i]].WriteMessage(1, res.JSON())
	// }
}
