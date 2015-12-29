package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/syhlion/gwspack"
)

func WSConnect(w http.ResponseWriter, r *http.Request) {

	log.Info(r.RemoteAddr, " handshake start")
	params := mux.Vars(r)
	app_key := params["app_key"]
	user_tag := params["user_tag"]
	if app_key == "" || user_tag == "" {
		log.Warn(r.RemoteAddr, " app_key & user_tag empty")
		http.Error(w, "app_key || user_tag empty", 404)
		return
	}
	app := gwspack.Get(app_key)
	println("test")
	c, err := app.Register(user_tag, w, r, nil)
	if err != nil {
		log.Warn(r.RemoteAddr, " ", err)
		http.Error(w, err.Error(), 403)
		return
	}

	log.Info(r.RemoteAddr, " login ", app_key, " Scuess")
	c.Listen()
	defer log.Info(r.RemoteAddr, " ", user_tag, " logout", " ", app_key)
	return
}
