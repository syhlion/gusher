package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func UnregisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	}
}

func ListClientHandler(w http.ResponseWriter, r *http.Request) {

}

type NilResult struct {
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
	request_ip := r.RemoteAddr

	if app_name == "" || request_ip == "" {
		log.Warn(r.RemoteAddr, " ", "app_name || request_ip empty")
		w.WriteHeader(404)
		result := NilResult{Message: "app_name || request_ip empty"}
		json.NewEncoder(w).Encode(result)
		return
	}
	app_key, err := appdata.Register(app_name, request_ip)

	if err != nil {
		log.Warn(r.RemoteAddr, " ", err)
		result := NilResult{Message: "Insert Error"}
		json.NewEncoder(w).Encode(result)
		return
	}

	result := AppResult{
		AppName:   app_name,
		AppKey:    app_key,
		RequestIP: request_ip,
	}
	json.NewEncoder(w).Encode(result)
}

func ListAppHandler(w http.ResponseWriter, r *http.Request) {

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
		nr := NilResult{Message: "app_key || content empty"}
		w.WriteHeader(400)
		log.Warn(r.RemoteAddr, " empty app_key || content")
		json.NewEncoder(w).Encode(nr)
		return
	}

	app, err := collection.Get(app_key)
	if err != nil {
		nr := NilResult{Message: err.Error()}
		log.Warn(r.RemoteAddr, " ", app_key, " ", err)
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(nr)
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
