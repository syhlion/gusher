package handle

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/syhlion/gopusher/model"
	"github.com/syhlion/gopusher/module/log"
	"net/http"
	"strconv"
)

func (h *Handler) AppList(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		log.Logger.Warn("ParseError")
	}
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		log.Logger.Warn("ParseError")
	}

	rs, err := h.AppData.GetAll()
	if err != nil {
		log.Logger.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	//pagination
	offset := (page - 1) * limit
	count := 0
	var tmprs []model.AppDataResult
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

	log.Logger.Info(r.RemoteAddr, " ListApp Scuess")
	json.NewEncoder(w).Encode(result)

}

func (h *Handler) Unregister(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app_key := params["app_key"]
	if app_key == "" {
		log.Logger.Warn(r.RemoteAddr, " app_key empty")
		http.Error(w, "app_key empty", 404)
		return
	}
	err := h.AppData.Delete(app_key)
	if err != nil {

		log.Logger.Warn(r.RemoteAddr, " ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	nr := NormalResult{Message: "Scuess"}
	log.Logger.Info(r.RemoteAddr, " ", app_key, " Unregister Scuess")
	json.NewEncoder(w).Encode(nr)

}

func (h *Handler) ListClient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	app_key := params["app_key"]
	if app_key == "" {
		log.Logger.Warn(r.RemoteAddr, " app_key empty")
		http.Error(w, "app_key empty", 404)
		return
	}
	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		log.Logger.Warn("ParseError")
	}
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		log.Logger.Warn("ParseError")
	}

	app, err := h.Collection.Get(app_key)

	if err != nil {
		log.Logger.Warn(r.RemoteAddr, " ", app_key, " ", err)
		http.Error(w, err.Error(), 403)
		return
	}
	onlineUsers := app.GetAllUserTag()

	//pagination
	offset := (page - 1) * limit
	count := 0
	var tmprs []string
	for n, v := range onlineUsers {
		if n >= offset {
			count++
			tmprs = append(tmprs, v)
			if count == limit {
				break

			}

		}
	}
	lo := ListOnlineResult{
		AppKey:   app_key,
		Total:    len(onlineUsers),
		UserTags: tmprs,
		Limit:    limit,
		Page:     page,
	}
	log.Logger.Info(r.RemoteAddr, " GetAppUsers")
	json.NewEncoder(w).Encode(lo)

}

//註冊
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	app_name := r.FormValue("app_name")
	auth_password := r.FormValue("auth_password")
	auth_account := r.FormValue("auth_account")
	request_ip := r.RemoteAddr

	if app_name == "" || request_ip == "" || auth_password == "" || auth_account == "" {
		log.Logger.Warn(r.RemoteAddr, " ", "app_name || request_ip empty")
		http.Error(w, "app_name || request_op empty", 404)
		return
	}
	app_key, err := h.AppData.Register(app_name, auth_account, auth_password, request_ip)

	if err != nil {
		log.Logger.Warn(r.RemoteAddr, " ", err)
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

//Push
func (h *Handler) Push(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	app_key := params["app_key"]
	content := r.FormValue("content")

	//option
	user_tag := r.FormValue("user_tag")

	if app_key == "" || content == "" {
		log.Logger.Warn(r.RemoteAddr, " empty app_key || content")
		http.Error(w, "app_key || content empty", 400)
		return
	}

	app, err := h.Collection.Get(app_key)
	if err != nil {
		log.Logger.Warn(r.RemoteAddr, " ", app_key, " ", err)
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

	log.Logger.Info(r.RemoteAddr, " message send ", content)
	json.NewEncoder(w).Encode(pushResult)
}
