package main

import (
	"github.com/gorilla/mux"
)

func Router() (router *mux.Router) {
	router = mux.NewRouter()

	// ws handshake
	router.HandleFunc("/ws/{app_key}/{user_tag}", WSHandler).Methods("GET")

	//push message api
	router.HandleFunc("/api/push/{app_key}", PushHandler).Methods("POST")

	//register user
	router.HandleFunc("/api/register", RegisterHandler).Methods("POST")

	//unregister
	//router.HandleFunc("/unregister", UnregisterHandler).Methods("POST")

	//list how many client
	router.HandleFunc("/api/{app_key}/listonlineuser", ListClientHandler).Methods("GET")
	return
}
