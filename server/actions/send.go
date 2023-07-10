package actions

import (
	"dicedasher/st"
	"dicedasher/storage"

	"golang.org/x/exp/slices"
)

func auth_error(response st.Response) { // Response with AUTH ERROR message
	response.Status = "auth_error"
	response.Data = ""
	storage.Clients[response.Player_id].WriteMessage(1, response.JSON())
}

func Send(response st.Response) {
	room, ok := storage.RoomStorage[response.Room_id]
	if !ok { // No room error 
		auth_error(response)
		return 
	}
	players := room.Players
	if !slices.Contains(players, response.Player_id) { // No player error 
		auth_error(response)
		return
	}
	for i := 0; i < len(players); i++ { // Send message to every client in the room 
		_, exists := storage.Clients[players[i]]
		if exists {
			storage.Clients[players[i]].WriteMessage(1, response.JSON())
		}
	}
}