package actions

import (
	"dicedasher/st"
	"dicedasher/storage"
)

func StartGame(request st.Request) bool {
	player_id := request.Player_id
	room := storage.RoomStorage[request.Room_id]
	if room.Master == player_id {
		room.IsOpened = false 
		storage.RoomStorage[request.Room_id] = room 
		Send(st.Response{
			Room_id: request.Room_id,
			Player_id: player_id,
			Action: request.Action,
			Status: "success",
			Data: "",
		})
		return true 
	} 
	return false
}