package storage

import (
	"dicedasher/st"

	"github.com/gorilla/websocket"
)

var RoomStorage = make(map[string]st.Room) 

var Clients = make(map[string]*websocket.Conn)

var Users = make(map[string]st.User)

var AccessTokens = make(map[string]string) // access_token : player_id