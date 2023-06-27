package main

import (
	"io"
	"net/http"

	rooms "dicedasher/rooms"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `hello`)
	})

	r.HandleFunc("/newroom", rooms.NewRoom)

	http.ListenAndServe(":80", r)
}

/*

import (
    "math/rand"
    "time"
)

var rooms = make(map[string]*Room)

func CreateRoom () {
	newID := GenerateUniqueId()

	newRoom := &Room {
		ID = newID
		master = userID
}

func GenerateUniqueID () {
	id := ""
	seed := rand.NewSource(time.Now().UnixNano())
	timeRand := rand.New(seed)
	for i := 0; i < 10; i++ {
		id += string(timeRand.Intn(58) + 65)
	}
}

*/
