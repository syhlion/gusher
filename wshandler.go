package main

import (
	"encoding/json"
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
		nr := NilResult{Message: "app_key & user_tag empty"}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(nr)
		return
	}

	//collection join
	app, err := collection.Join(app_key)
	log.Debug("test")
	if err != nil {
		log.Warn(r.RemoteAddr, " ", app_key, " ", err)
		nr := NilResult{Message: err.Error()}
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(nr)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn(r.RemoteAddr, " ", err)
		nr := NilResult{Message: err.Error()}
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(nr)
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
