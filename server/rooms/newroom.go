package rooms

import (
	"dicedasher/st"
	"dicedasher/storage"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func generateID() string {
	lib := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	id := ""
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	for i := 0; i < 10; i++ {
		id += string(lib[(timeRand.Intn(len(lib)))])
	}
	return id
}

func NewRoom(c *gin.Context) {

	var room st.Room = st.Room{
		ID: generateID(),
		Master: generateID(),
		IsOpened: true,
	} // Create Room structure
	room.Players = append(room.Players, room.Master) // Append master id to players array

	storage.RoomStorage[room.ID] = room // Store Room structure in RoomStorage

	res := &st.Response{ // HTTP Response structure
		Room_id: room.ID, 
		Player_id: room.Master, 
		Action: "create_room",
		Status: "success", 
		Data: "{}",
	}
	c.JSON(200, res) // Send HTTP Response
}
