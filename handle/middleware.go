package handle

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/syhlion/gopusher/model"
	"github.com/syhlion/gopusher/module/config"
	"github.com/syhlion/gopusher/module/log"
	"net/http"
	"strings"
)

type Middleware struct {
	AppData *model.AppData
	Config  *config.Config
}

//Middleware use
func (m *Middleware) Use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func (m *Middleware) AppKeyVerity(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		if !m.AppData.IsExist(params["app_key"]) {
			log.Logger.Warn(r.RemoteAddr, " ", params["app_key"]+" app_key does not exist")
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
			log.Logger.Warn(r.RemoteAddr, "auth Error")
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
			log.Logger.Warn(r.RemoteAddr, " auth param empty")
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
				log.Logger.Warn(r.RemoteAddr, " ", err.Error())
				http.Error(w, err.Error(), 401)
				return
			}
			account = data.AuthAccount
			password = data.AuthPassword
		}

		if pair[0] != account || pair[1] != password {
			log.Logger.Warn(r.RemoteAddr, " auth error "+pair[0]+" "+pair[1])
			http.Error(w, "Not authorized", 401)
			return
		}
		h.ServeHTTP(w, r)

	}
}
