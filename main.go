package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

var collection = NewCollection()

var log = logrus.New()

func main() {
	// log init

	go collection.run()
	r := mux.NewRouter()
	r.HandleFunc("/ws", WSHandler).Methods("GET")
	r.HandleFunc("/push", PushHandler)
	log.Info("Server Start")
	http.ListenAndServe(":8001", r)
}
