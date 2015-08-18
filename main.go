package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

var collection = NewCollection()

var log = logrus.New()

func logHttpRequestInfo(r *http.Request) {
	log.Info()
}
func logHttpRequestWarn(r *http.Request) {
	log.Warn()
}
func logWsInfo(c *Client) {
	log.Info()
}

func logWsWarn(c *Client) {
	log.Info()
}

func main() {
	// log init

	go collection.run()
	r := mux.NewRouter()

	// ws handshake
	r.HandleFunc("/ws/{key}", WSHandler).Methods("GET")

	//push message api
	r.HandleFunc("/push/{key}", PushHandler).Methods("POST")

	//register user
	r.HandleFunc("/register", RegisterHandler).Methods("POST")

	//unregister
	r.HandleFunc("/unregister", UnregisterHandler).Methods("POST")

	//list how many client
	r.HandleFunc("/{key}/list", ListClientHandler).Methods("GET")
	log.Info("Server Start")
	http.ListenAndServe(":8001", r)
}
