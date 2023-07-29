package actions

import (
	"github.com/gorilla/websocket"
	"time"
)

func TryCloseNormally(wsConn *websocket.Conn) error {
	closeNormalClosure := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := wsConn.WriteControl(websocket.CloseMessage, closeNormalClosure, time.Now().Add(time.Second)); err != nil {
		return err
	}
	return wsConn.Close()
}
