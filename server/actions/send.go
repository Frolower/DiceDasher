package actions

import (
	"dicedasher/st"
	"dicedasher/storage"
)

func send(response st.Response) {
	players := storage.RoomStorage[response.Room_id].Players // TODO: Handle error here
	for i := 0; i < len(players); i++ {
		storage.Clients[players[i]].WriteMessage(1, response.JSON())
	}
}