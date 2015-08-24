package router

import (
	"github.com/gorilla/mux"
	"github.com/syhlion/gopusher/handler"
)

func Router() (router *mux.Router) {
	router = mux.NewRouter()

	// ws handshake
	router.HandleFunc("/ws/{app_key}/{user_tag}", handler.WSHandler).Methods("GET")

	//push message api
	router.HandleFunc("/api/push/{app_key}", handler.PushHandler).Methods("POST")

	//register user
	router.HandleFunc("/api/register", handler.RegisterHandler).Methods("POST")

	//unregister
	router.HandleFunc("/api/{app_key}/unregister", handler.UnregisterHandler).Methods("DELETE")

	//list app
	router.HandleFunc("/api/app-list", handler.AppListHandler).Methods("GET")

	//list how many client
	router.HandleFunc("/api/{app_key}/listonlineuser", handler.ListClientHandler).Methods("GET")
	return
}
