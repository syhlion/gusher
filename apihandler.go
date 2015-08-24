package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type AppDataResult struct {
	AppKey       string `json:"app_key"`
	AppName      string `json:"app_name"`
	AuthAccount  string `json:"auth_account"`
	AuthPassword string `json:"auth_password"`
	RequestIP    string `json:"request_ip"`
	Date         string `json:"date"`
	Timestamp    string `json:"timestamp"`
}

func AppListHandler(w http.ResponseWriter, r *http.Request) {
	rs, err := appdata.GetAll()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Info(r.RemoteAddr, " ListApp Scuess")
	json.NewEncoder(w).Encode(rs)

}

func UnregisterHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app_key := params["app_key"]
	if app_key == "" {
		log.Warn(r.RemoteAddr, " app_key empty")
		http.Error(w, "app_key empty", 404)
		return
	}
	err := appdata.Delete(app_key)
	if err != nil {

		log.Warn(r.RemoteAddr, " ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	nr := NormalResult{Message: "Scuess"}
	log.Info(r.RemoteAddr, " ", app_key, " Unregister Scuess")
	json.NewEncoder(w).Encode(nr)

}

type ListOnlineResult struct {
	AppKey          string   `json:"app_key"`
	TotalOnlineUser int      `json:"total_online_user"`
	OnlineUser      []string `json:"online_user"`
}

func ListClientHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app_key := params["app_key"]
	if app_key == "" {
		log.Warn(r.RemoteAddr, " app_key empty")
		http.Error(w, "app_key empty", 404)
		return
	}

	app, err := collection.Get(app_key)

	if err != nil {
		log.Warn(r.RemoteAddr, " ", app_key, " ", err)
		http.Error(w, err.Error(), 403)
		return
	}
	onlineUsers := app.GetAllUserTag()

	lo := ListOnlineResult{
		AppKey:          app_key,
		TotalOnlineUser: len(onlineUsers),
		OnlineUser:      onlineUsers,
	}
	log.Info(r.RemoteAddr, " GetAppUsers")
	json.NewEncoder(w).Encode(lo)

}

type NormalResult struct {
	Message string `json:"message"`
}

type AppResult struct {
	AppName   string `json:"app_name"`
	AppKey    string `json:"app_key"`
	RequestIP string `json:"request_ip"`
}

//註冊
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	app_name := r.FormValue("app_name")
	auth_password := r.FormValue("auth_password")
	auth_account := r.FormValue("auth_account")
	request_ip := r.RemoteAddr

	if app_name == "" || request_ip == "" || auth_password == "" || auth_account == "" {
		log.Warn(r.RemoteAddr, " ", "app_name || request_ip empty")
		http.Error(w, "app_name || request_op empty", 404)
		return
	}
	app_key, err := appdata.Register(app_name, auth_account, auth_password, request_ip)

	if err != nil {
		log.Warn(r.RemoteAddr, " ", err)
		http.Error(w, "Insert Error", 500)
		return
	}

	result := AppResult{
		AppName:   app_name,
		AppKey:    app_key,
		RequestIP: request_ip,
	}
	json.NewEncoder(w).Encode(result)
}

type PushResult struct {
	AppKey  string `json:"app_key"`
	Content string `json:"content"`
	UserTag string `json:"user_tag"`
	Total   int    `json:"total"`
}

//Push
func PushHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	app_key := params["app_key"]
	content := r.FormValue("content")

	//option
	user_tag := r.FormValue("user_tag")

	if app_key == "" || content == "" {
		log.Warn(r.RemoteAddr, " empty app_key || content")
		http.Error(w, "app_key || content empty", 400)
		return
	}

	app, err := collection.Get(app_key)
	if err != nil {
		log.Warn(r.RemoteAddr, " ", app_key, " ", err)
		http.Error(w, err.Error(), 403)
		return
	}
	totalResult := 0
	b := []byte(content)
	if user_tag == "" {
		app.Boradcast <- b
		totalResult = len(app.Connections)
	} else {
		m := make(map[string][]byte)
		m[user_tag] = b
		app.Assign <- m
		totalResult = <-app.AssignTotalResult
	}

	pushResult := PushResult{
		AppKey:  app_key,
		Content: content,
		UserTag: user_tag,
		Total:   totalResult,
	}

	log.Info(r.RemoteAddr, " message send ", content)
	json.NewEncoder(w).Encode(pushResult)
}
