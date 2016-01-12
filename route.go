package main

import (
	"github.com/gorilla/mux"
)

func PublicRouter() (router *mux.Router) {
	router = mux.NewRouter()

	//middleware.use()//先寫後執行
	// ws handshake
	router.HandleFunc("/ws/{app_key}/{user_tag}", MiddlewareUse(WSConnect, ConnectWebHook, AppKeyVerity, LogHttpRequest)).Methods("GET")

	return
}
func PrivateRouter() (router *mux.Router) {
	router = mux.NewRouter()

	//middleware.use()//先寫後執行
	// ws handshake
	//push message api
	router.HandleFunc("/api/push/{app_key}", MiddlewareUse(Push, AppKeyVerity, BasicAuth, LogHttpRequest)).Methods("POST")

	//register user
	router.HandleFunc("/api/register", MiddlewareUse(Register, BasicAuth, LogHttpRequest)).Methods("POST")

	//unregister
	router.HandleFunc("/api/{app_key}/unregister", MiddlewareUse(Unregister, AppKeyVerity, BasicAuth, LogHttpRequest)).Methods("DELETE")

	//list app
	router.HandleFunc("/api/app-list", MiddlewareUse(AppList, BasicAuth, LogHttpRequest)).Methods("GET")

	//list how many client
	router.HandleFunc("/api/{app_key}/listonlineuser", MiddlewareUse(ListClient, AppKeyVerity, BasicAuth, LogHttpRequest)).Methods("GET")
	return
}
