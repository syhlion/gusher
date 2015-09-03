package handle

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/syhlion/gusher/model"
	"github.com/syhlion/gusher/module/config"
	"github.com/syhlion/gusher/module/requestworker"
)

type Middleware struct {
	AppData *model.AppData
	Config  *config.Config
	Worker  *requestworker.Worker
}

//Middleware use
func (m *Middleware) Use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func (m *Middleware) LogHttpRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RemoteAddr, " ", r.Method, " ", r.RequestURI, " ", r.Header.Get("Authorization"))
		h.ServeHTTP(w, r)

	}
}

func (m *Middleware) ConnectWebHook(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		token := r.FormValue("token")
		data, err := m.AppData.Get(params["app_key"])
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}

		// no fill in connect_hook
		hook_url := data.ConnectHook
		if hook_url == "" {
			h.ServeHTTP(w, r)
			return
		}

		//fill in connect_hook bug url parse error
		u, err := url.Parse(hook_url)
		if err != nil {
			log.Warn(r.RemoteAddr, " ", params["app_key"], " ", err.Error())
			http.Error(w, "hook url error", 404)
			return
		}

		//hook  url requset
		v := url.Values{}
		v.Add("token", token)
		req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(v.Encode()))
		if err != nil {
			log.Warn(r.RemoteAddr, " ", err.Error())
			http.Error(w, err.Error(), 404)
			return
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(v.Encode())))
		job := &requestworker.Job{
			Resq:   req,
			Result: make(chan requestworker.Result),
		}
		m.Worker.JobQuene <- job
		rs := <-job.Result
		if rs.Err != nil {
			log.Warn(r.RemoteAddr, " ", rs.Err.Error())
			http.Error(w, rs.Err.Error(), 404)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func (m *Middleware) AppKeyVerity(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		if !m.AppData.IsExist(params["app_key"]) {
			log.Warn(r.RemoteAddr, " ", params["app_key"]+" app_key does not exist")
			http.Error(w, "app_key does not exist", 404)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func (m *Middleware) BasicAuth(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenicate", `Basic realm="Restricted`)
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			log.Warn(r.RemoteAddr, "auth Error")
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			log.Warn(r.RemoteAddr, " auth param empty")
			http.Error(w, "Not authorized", 401)
			return
		}

		params := mux.Vars(r)
		var account string
		var password string

		//super admin 可通過任何api
		if pair[0] == m.Config.AuthAccount && pair[1] == m.Config.AuthPassword {
			account = m.Config.AuthAccount
			password = m.Config.AuthPassword
			h.ServeHTTP(w, r)
			return
		}

		if params["app_key"] != "" {
			data, err := m.AppData.Get(params["app_key"])
			if err != nil {
				log.Warn(r.RemoteAddr, " ", err.Error())
				http.Error(w, err.Error(), 401)
				return
			}
			account = data.AuthAccount
			password = data.AuthPassword
		}

		if pair[0] != account || pair[1] != password {
			log.Warn(r.RemoteAddr, " auth error "+pair[0]+" "+pair[1])
			http.Error(w, "Not authorized", 401)
			return
		}
		h.ServeHTTP(w, r)

	}
}
