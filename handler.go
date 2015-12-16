package main

import (
	"github.com/gorilla/mux"
	"github.com/syhlion/requestwork"
)

type Handler struct {
}

func (h *Handler) Router(appdata *model.AppData, collection *core.Collection, config *config.Config, worker *requestwork.Worker) (router *mux.Router) {
	router = mux.NewRouter()

	handler := &handle.Handler{appdata, collection}
	middleware := &handle.Middleware{appdata, config, worker}

	//middleware.use()//先寫後執行

	// ws handshake
	router.HandleFunc("/ws/{app_key}/{user_tag}", middleware.Use(handler.WS, middleware.ConnectWebHook, middleware.AppKeyVerity, middleware.LogHttpRequest)).Methods("GET")

	//push message api
	router.HandleFunc("/api/push/{app_key}", middleware.Use(handler.Push, middleware.AppKeyVerity, middleware.BasicAuth, middleware.LogHttpRequest)).Methods("POST")

	//register user
	router.HandleFunc("/api/register", middleware.Use(handler.Register, middleware.BasicAuth, middleware.AllowAccessApiIP, middleware.LogHttpRequest)).Methods("POST")

	//unregister
	router.HandleFunc("/api/{app_key}/unregister", middleware.Use(handler.Unregister, middleware.AppKeyVerity, middleware.BasicAuth, middleware.AllowAccessApiIP, middleware.LogHttpRequest)).Methods("DELETE")

	//list app
	router.HandleFunc("/api/app-list/{limit:[0-9]+}/{page:[0-9]+}", middleware.Use(handler.AppList, middleware.BasicAuth, middleware.AllowAccessApiIP, middleware.LogHttpRequest)).Methods("GET")

	//list how many client
	router.HandleFunc("/api/{app_key}/listonlineuser/{limit:[0-9]+}/{page:[0-9]+}", middleware.Use(handler.ListClient, middleware.AppKeyVerity, middleware.BasicAuth, middleware.AllowAccessApiIP, middleware.LogHttpRequest)).Methods("GET")
	return
}
