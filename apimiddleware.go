package main

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

//Middleware use
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func BasicAuth(h http.HandlerFunc) http.HandlerFunc {

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

		//由 router是否有 app_key 來判斷是否是 super admin
		params := mux.Vars(r)
		var account string
		var password string
		if params["app_key"] == "" {

		} else {
			data, err := appdata.Get(params["app_key"])
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
