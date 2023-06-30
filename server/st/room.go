package st

type Room struct {
	ID      string
	Master  string
	Players []string
	IsOpened bool
}

func (g Room) isPlayerConnected(player_id string) bool {
	for i:=0; i<len(g.Players); i++ {
		if g.Players[i] == player_id {
			return true 
		}
	}
	return false 
}