package actions

import (
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
	r := RollDice(100) // Generate values 
	response := st.Response{ // Create response message
			Room_id: request.Room_id,
			Player_id: request.Player_id,
			Action: request.Action,
	}
	data, _ := json.Marshal(r) // Turn values into string
	response.Data = string(data)
	send(response) // Send message to every player in the room
}
