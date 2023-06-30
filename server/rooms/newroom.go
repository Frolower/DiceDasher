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
	var room st.Room
	room.ID = generateID()
	room.Master = generateID()
	room.Players = append(room.Players, room.Master)
	room.IsOpened = true

	storage.RoomStorage[room.ID] = room
	res := &st.Response{Room_id: room.ID, Player_id: room.Master, Action: "create_room", Data: make(map[string]string)}
	c.JSON(200, res)
}
