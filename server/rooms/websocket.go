package rooms

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn) // Map with clients {player_id : connection}

var upgrader = websocket.Upgrader{ // Upgrade HTTP connection to WebSocket
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
