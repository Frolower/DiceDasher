package rooms

type Room struct {
	ID      string
	master  string
	players []string
	isOpened bool
}

func (this Room) isPlayerConnected(player_id string) bool {
	for i:=0; i<len(this.players); i++ {
		if this.players[i] == player_id {
			return true 
		}
	}
	return false 
}