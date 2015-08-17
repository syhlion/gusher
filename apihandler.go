package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("post")
	if r.Method != "POST" {
		log.Fatal("not post")
		return
	}

	key := r.FormValue("key")
	content := r.FormValue("content")

	log.Println(key, content)
	if key == "" || content == "" {
		log.Fatal("no content")
		return
	}

	log.Println("read get")
	room, _ := collection.Get(key)
	log.Println("get OK")
	if room == nil {
		log.Fatal("no user")
		return
	}
	room.Boradcast <- []byte(content)

	fmt.Fprintf(w, "Scuess, %q", html.EscapeString(r.URL.Path))
}
