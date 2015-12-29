package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/syhlion/gwspack"
)

func AppList(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		log.Warn("ParseError")
	}
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		log.Warn("ParseError")
	}

	rs, err := Model.GetAll()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	//pagination
	offset := (page - 1) * limit
	count := 0
	var tmprs []AppDataResult
	for n, v := range rs {
		if n >= offset {
			count++
			tmprs = append(tmprs, v)
			if count == limit {
				break

			}

		}
	}
	result := AppListResult{
		Limit: limit,
		Page:  page,
		Total: len(rs),
		Data:  tmprs,
	}

	log.Info(r.RemoteAddr, " ListApp Scuess")
	json.NewEncoder(w).Encode(result)

}

func Unregister(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app_key := params["app_key"]
	if app_key == "" {
		log.Warn(r.RemoteAddr, " app_key empty")
		http.Error(w, "app_key empty", 404)
		return
	}
	err := Model.Delete(app_key)
	if err != nil {

		log.Warn(r.RemoteAddr, " ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	nr := NormalResult{Message: "Scuess"}
	log.Info(r.RemoteAddr, " ", app_key, " Unregister Scuess")
	json.NewEncoder(w).Encode(nr)

}

func ListClient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app_key := params["app_key"]
	if app_key == "" {
		log.Warn(r.RemoteAddr, " app_key empty")
		http.Error(w, "app_key empty", 404)
		return
	}
	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		log.Warn("ParseError")
	}
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		log.Warn("ParseError")
	}

	app := gwspack.Get(app_key)

	onlineUsers := app.List()

	//pagination
	/*
		offset := (page - 1) * limit
		count := 0
		var tmprs []string
		n:=0
		for k, _ := range onlineUsers {
			n++
			if n >= offset {
				count++
				tmprs = append(tmprs, k)
				if count == limit {
					break

				}

			}
		}*/
	lo := ListOnlineResult{
		AppKey: app_key,
		Total:  len(onlineUsers),
		Limit:  limit,
		Page:   page,
	}
	log.Info(r.RemoteAddr, " GetAppUsers")
	json.NewEncoder(w).Encode(lo)

}

//註冊
func Register(w http.ResponseWriter, r *http.Request) {
	app_key := r.FormValue("app_key")
	auth_password := r.FormValue("auth_password")
	auth_account := r.FormValue("auth_account")
	connect_hook := r.FormValue("connect_hook")
	request_ip := r.RemoteAddr

	if app_key == "" || request_ip == "" || auth_password == "" || auth_account == "" {
		log.Warn(r.RemoteAddr, " ", "app_name || request_ip empty")
		http.Error(w, "app_name || request_op empty", 404)
		return
	}

	//bcrypt encoding
	hash_password, err := bcrypt.GenerateFromPassword([]byte(auth_account+auth_password), 5)
	if err != nil {
		log.Warn(r.RemoteAddr, " ", app_key, " ", auth_password, " hash error")
		http.Error(w, "hash error", 404)
		return
	}
	auth_password = string(hash_password)

	err = Model.Register(app_key, auth_account, auth_password, connect_hook, request_ip)

	if err != nil {
		log.Warn(r.RemoteAddr, " ", err)
		http.Error(w, "app_key repeat", 500)
		return
	}

	result := AppResult{
		AppKey:      app_key,
		ConnectHook: connect_hook,
		RequestIP:   request_ip,
	}
	json.NewEncoder(w).Encode(result)
}

//Push
func Push(w http.ResponseWriter, r *http.Request) {

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

	app := gwspack.Get(app_key)
	b := []byte(content)
	if user_tag == "" {
		app.SendAll(b)
	} else {
		app.SendByRegex(user_tag, b)

	}

	pushResult := &PushResult{
		AppKey:  app_key,
		Content: content,
		UserTag: user_tag,
	}

	log.Info(r.RemoteAddr, " message send ", content)
	json.NewEncoder(w).Encode(pushResult)
}
