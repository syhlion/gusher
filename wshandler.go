package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WSHandler(w http.ResponseWriter, r *http.Request) {

	log.Info(r.RemoteAddr, " handshake start")
	params := mux.Vars(r)

	app_key := params["app_key"]
	user_tag := params["user_tag"]
	if app_key == "" || user_tag == "" {
		log.Warn(r.RemoteAddr, " app_key & user_tag empty")
		http.Error(w, "app_key || user_tag empty", 404)
		return
	}

	//collection join
	app, err := collection.Join(app_key)
	if err != nil {
		log.Warn(r.RemoteAddr, " ", app_key, " ", err)
		http.Error(w, err.Error(), 403)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn(r.RemoteAddr, " ", err)
		http.Error(w, err.Error(), 403)
		return
	}

	// new client
	client := NewClient(user_tag, ws, app)

	// register client
	app.Register <- client
	log.Info(r.RemoteAddr, " login ", app_key, " Scuess")
	go client.WritePump()
	client.ReadPump()
	defer log.Info(r.RemoteAddr, " logout")
}
