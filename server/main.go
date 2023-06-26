package main

import (
	"io"
	"net/http"

	rooms "./rooms"
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