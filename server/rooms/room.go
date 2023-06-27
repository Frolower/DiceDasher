package rooms

type Room struct {
	ID      string
	master  string
	players []string
	isOpened bool
}