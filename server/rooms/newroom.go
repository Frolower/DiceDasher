package rooms

import (
	"dicedasher/auth"
	"dicedasher/st"
	"dicedasher/storage"
	"fmt"
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
	token, err := c.Cookie("access_token")
	fmt.Println(token)
	if err != nil {
		fmt.Println(err)
		c.JSON(401, nil)
		return
	}
	user_id := auth.Auth(token)

	if user_id == "" {
		c.JSON(401, nil) 
		return
	}

	var room st.Room = st.Room{
		ID: generateID(),
		Master: user_id,
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
