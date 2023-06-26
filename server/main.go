package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `hello`)
	})
	http.ListenAndServe(":80", r)
}