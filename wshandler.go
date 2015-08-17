package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Metho not allowed", 405)
		//TODO 補log
		log.Println("error")
		return
	}

	u := r.FormValue("key")
	if u == "" {
		//TODO 補log
		return
	}
	log.Println(u)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//TODO 補logo
		return
	}
	coll, err := collection.Join(u)
	if err != nil {
		return
	}

	client := &Client{
		ws:   ws,
		room: coll,
		send: make(chan []byte),
	}

	coll.Register <- client
	go client.writePump()
	client.readPump()
}
