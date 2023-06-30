package storage

import (
	"dicedasher/st"

	"github.com/gorilla/websocket"
)

var RoomStorage = make(map[string]st.Room) 

var Clients = make(map[string]*websocket.Conn)