package main

import (
	"github.com/gorilla/mux"
)

func Router() (router *mux.Router) {
	router = mux.NewRouter()

	/****************** client use *************************/

	// ws handshake
	router.HandleFunc("/ws/{app_key}/{user_tag}", WSHandler).Methods("GET")

	/****************** Super admin use *************************/

	//register user
	router.HandleFunc("/api/register", RegisterHandler).Methods("POST")

	//unregister
	router.HandleFunc("/api/{app_key}/unregister", UnregisterHandler).Methods("DELETE")

	//list app
	router.HandleFunc("/api/app-list", AppListHandler).Methods("GET")

	/****************** admin use *************************/

	//push message api
	router.HandleFunc("/api/push/{app_key}", PushHandler).Methods("POST")

	//list how many client
	router.HandleFunc("/api/{app_key}/listonlineuser", ListClientHandler).Methods("GET")
	return
}
