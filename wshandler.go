package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
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
		return
	}

	keys := r.FormValue("keys")
	if keys == "" {
		//TODO 補log
		return
	}

	keys_array := strings.SplitN(keys, ":", 2)
	keys_total := len(keys_array)
	if keys_total != 2 {
		return
	}
	app_key := keys_array[0]
	user_token := keys_array[1]

	//csrf
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//TODO 補logo
		return
	}

	//collection join
	coll, err := collection.Join(app_key)
	if err != nil {
		return
	}

	// new client
	client := &Client{
		token: user_token,
		ws:    ws,
		room:  coll,
		send:  make(chan []byte),
	}

	// register client
	coll.Register <- client
	client.writePump()
	//client.readPump()
}
