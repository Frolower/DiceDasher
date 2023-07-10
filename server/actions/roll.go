package actions

import (
	"dicedasher/st"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

func rollDice(size int) []int {
	var result []int
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	numIter := (timeRand.Intn(5) + 10)
	for i := 0; i < numIter; i++ {
	 result = append(result, timeRand.Intn(size)+1)
	}
	return result
   }

func Roll(request st.Request) st.Response {
	rollTypes := []string{"1d4", "1d6", "1d8", "1d10", "1d12", "1d20", "1d100"}
	diceType := request.Data.Get("dice").String()
	if !slices.Contains(rollTypes, diceType) { // Check for dice type
		response := st.Response{ // Create response message
			Room_id: request.Room_id,
			Player_id: request.Player_id,
			Action: request.Action,
			Status: "error",
			Data: "bad_request",
		}
		return response	
	}
	diceMax, _ := strconv.Atoi(strings.Split(diceType, "d")[1])
	r := rollDice(diceMax)
	 // Generate values 
	fmt.Println(r)
	response := st.Response{ // Create response message
			Room_id: request.Room_id,
			Player_id: request.Player_id,
			Action: request.Action,
			Status: "success",
	}
	data, _ := json.Marshal(r) // Turn values into string
	response.Data = string(data)
	Send(response) // Send message to every player in the room
	return response
}
