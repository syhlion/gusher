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
		return
	}

	//確認db 是否存在
	if !appdata.IsExist(app_key) {
		log.Warn(r.RemoteAddr, " ", app_key, " no exist")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn(r.RemoteAddr, " ", err)
		return
	}

	//collection join
	app := collection.Join(app_key)
	// new client
	client := NewClient(user_tag, ws, app)

	// register client
	app.Register <- client
	log.Info(r.RemoteAddr, " login ", app_key, " Scuess")
	go client.WritePump()
	client.ReadPump()
}
