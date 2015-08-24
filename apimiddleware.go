package main

import (
	"encoding/base64"
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
		if pair[0] != "scott" || pair[1] != "xxxx" {
			log.Warn(r.RemoteAddr, " auth error "+pair[0]+" "+pair[1])
			http.Error(w, "Not authorized", 401)
			return
		}
		h.ServeHTTP(w, r)

	}
}
