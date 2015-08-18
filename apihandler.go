package main

import (
	"fmt"
	"html"
	"net/http"
)

//創建
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		return
	}
}

//Push
func PushHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		return
	}

	key := r.FormValue("key")
	content := r.FormValue("content")

	if key == "" || content == "" {
		return
	}

	room, _ := collection.Get(key)
	if room == nil {
		return
	}
	room.Boradcast <- []byte(content)

	fmt.Fprintf(w, "Scuess, %q", html.EscapeString(r.URL.Path))
}
