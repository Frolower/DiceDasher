package rooms

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

var RoomStorage = make(map[string]Room) // Map with all rooms

func generateID() string {
	lib := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	id := ""
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	for i := 0; i < 100; i++ {
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

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	clients[room.master] = ws

	for { // Reading incoming messages
		res := &Response{Room_id: room.ID, Player_id: room.master} // Default info answer
		ws.WriteMessage(1, res.JSON())

		_, message, _ := ws.ReadMessage() // TODO: Handle error

		req := Request{}
		req.fromJSON(message)
		req.handleAction()

	}
}
