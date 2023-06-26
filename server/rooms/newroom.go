package rooms

import (
	"fmt"
	"net/http"
)

func NewRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Println("newroom")
}