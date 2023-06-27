package rooms

import (
	"fmt"
	"net/http"
	"math/rand"
	"time"
	"github.com/gin-gonic/gin"
)

var RoomStorage = make(map[string]Room)

func generateID () string {
	id := ""
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	for i := 0; i < 10; i++ {
		id += string(timeRand.Intn(58) + 65)
	}
	return id
}

func NewRoom(c *gin.Context) {
	var room Room
	room.ID = generateID()
	room.master = generateID()
	room.players = append(room.players, room.master)
	room.isOpened = true
	fmt.Println(room)

	RoomStorage[room.ID] = room
	fmt.Println(RoomStorage)

	c.JSON(http.StatusOK, gin.H{
		"room_id" : room.ID,
	})

}