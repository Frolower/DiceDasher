package actions

import (
	"dicedasher/st"
	"dicedasher/storage"
)

func Lock(request st.Request) bool {
	player_id := request.Player_id
	room := storage.RoomStorage[request.Room_id]
	if room.Master == player_id {
		
		if request.Action == "open_room" {
			room.IsOpened = true 
		} else if request.Action == "close_room" {
			room.IsOpened = false 
		} else {
			return false 
		}
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