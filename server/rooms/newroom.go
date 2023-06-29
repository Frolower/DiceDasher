package rooms

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

var RoomStorage = make(map[string]Room) // Map with all rooms

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
	var room Room
	room.ID = generateID()
	room.master = generateID()
	room.players = append(room.players, room.master)
	room.isOpened = true

	RoomStorage[room.ID] = room
	res := &Response{Room_id: room.ID, Player_id: room.master, Action: "create_room", Data: make(map[string]string)}
	c.JSON(200, res)
}
