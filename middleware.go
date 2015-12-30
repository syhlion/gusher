package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

//Middleware use
func MiddlewareUse(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func AllowAccessApiIP(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ips := strings.Split(r.RemoteAddr, ":")
		for _, allow := range GlobalConf.AllowAccessApiIP {
			if vailed, err := regexp.Compile(allow); err == nil {
				if vailed.MatchString(ips[0]) {
					h.ServeHTTP(w, r)
					return

				}
			}
		}
		log.Warn(r.RemoteAddr, " IP DENY")
		http.Error(w, "IP DENY", 404)
		return
	}
}

func LogHttpRequest(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RemoteAddr, " ", r.Method, " ", r.RequestURI, " ", r.Header.Get("Authorization"))
		h.ServeHTTP(w, r)

	}
}

func ConnectWebHook(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		token := r.FormValue("token")
		data, err := Model.Get(params["app_key"])
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
		resp, err := ReqWorker.Execute(req)
		if err != nil {
			log.Warn("Hook Error")
			http.Error(w, err.Error(), 404)
			return
		}
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Warn(err)
			http.Error(w, err.Error(), 404)
			return
		}
		ret := string(b)
		if ret != params["user_tag"] {
			log.Warn("Error user_tag " + params["user_tag"] + ", response user_tag " + ret)
			http.Error(w, "Error user_tag "+params["user_tag"]+", response user_tag "+ret, 404)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func AppKeyVerity(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		if !Model.IsExist(params["app_key"]) {
			log.Warn(r.RemoteAddr, " ", params["app_key"]+" app_key does not exist")
			http.Error(w, "app_key does not exist", 404)
			return
		}
		h.ServeHTTP(w, r)
	}
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

		params := mux.Vars(r)
		var account string
		var password string

		//super admin 可通過任何api
		if pair[0] == GlobalConf.AuthAccount && bcrypt.CompareHashAndPassword([]byte(GlobalConf.AuthPassword), []byte(pair[0]+pair[1])) == nil {
			h.ServeHTTP(w, r)
			return
		}

		if params["app_key"] != "" {
			data, err := Model.Get(params["app_key"])
			if err != nil {
				log.Warn(r.RemoteAddr, " ", err.Error())
				http.Error(w, err.Error(), 401)
				return
			}
			account = data.AuthAccount
			password = data.AuthPassword
		}

		if pair[0] != account || bcrypt.CompareHashAndPassword([]byte(password), []byte(pair[0]+pair[1])) != nil {
			log.Warn(r.RemoteAddr, " auth error "+pair[0]+" "+pair[1])
			http.Error(w, "Not authorized", 401)
			return
		}
		h.ServeHTTP(w, r)

	}
}
