package rooms

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ // Upgrade HTTP connection to WebSocket
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
