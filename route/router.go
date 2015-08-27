package route

import (
	"github.com/gorilla/mux"
	"github.com/syhlion/gopusher/core"
	"github.com/syhlion/gopusher/handle"
	"github.com/syhlion/gopusher/model"
	"github.com/syhlion/gopusher/module/config"
	"github.com/syhlion/gopusher/module/requestworker"
)

func Router(appdata *model.AppData, collection *core.Collection, config *config.Config, worker *requestworker.Worker) (router *mux.Router) {
	router = mux.NewRouter()

	handler := &handle.Handler{appdata, collection}
	middleware := &handle.Middleware{appdata, config, worker}

	// ws handshake
	router.HandleFunc("/ws/{app_key}/{user_tag}", middleware.Use(handler.WS, middleware.AppKeyVerity, middleware.ConnectWebHook)).Methods("GET")

	//push message api
	router.HandleFunc("/api/push/{app_key}", middleware.Use(handler.Push, middleware.AppKeyVerity, middleware.BasicAuth)).Methods("POST")

	//register user
	router.HandleFunc("/api/register", middleware.Use(handler.Register, middleware.BasicAuth)).Methods("POST")

	//unregister
	router.HandleFunc("/api/{app_key}/unregister", middleware.Use(handler.Unregister, middleware.AppKeyVerity, middleware.BasicAuth)).Methods("DELETE")

	//list app
	router.HandleFunc("/api/app-list/{limit:[0-9]}/{page:[0-9]}", middleware.Use(handler.AppList, middleware.BasicAuth)).Methods("GET")

	//list how many client
	router.HandleFunc("/api/{app_key}/listonlineuser/{limit:[0-9]}/{page:[0-9]}", middleware.Use(handler.ListClient, middleware.AppKeyVerity, middleware.BasicAuth)).Methods("GET")
	return
}
