package st

import (
	"math/rand"
	"time"
)

type User struct {
	Username string 
	Password_hash string 
	Email string 
	Id string 
}

func (g *User) GenerateID() {
	lib := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	id := ""
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	for i := 0; i < 10; i++ {
		id += string(lib[(timeRand.Intn(len(lib)))])
	}
	g.Id = id
}