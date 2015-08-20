package main

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"time"
)

var collection = NewCollection()

var log = logrus.New()

var appdata *AppData

func makeTimestamp() (t int64) {
	t = time.Now().UnixNano() / int64(time.Millisecond)
	return
}

func main() {
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}
	db, err := sql.Open("sqlite3", "./appdata.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table if not exists appdata (app_name,request_ip,app_key PRIMARY KEY,timestamp,date)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	appdata = NewAppData(db)
	go collection.run()
	r := mux.NewRouter()

	// ws handshake
	r.HandleFunc("/ws/{app_key}/{user_tag}", WSHandler).Methods("GET")

	//push message api
	r.HandleFunc("/push/{app_key}", PushHandler).Methods("POST")

	//register user
	r.HandleFunc("/register", RegisterHandler).Methods("POST")

	//unregister
	r.HandleFunc("/unregister", UnregisterHandler).Methods("POST")

	//list how many client
	r.HandleFunc("/{key}/listonlineuser", ListClientHandler).Methods("GET")

	r.HandleFunc("/listapp", ListAppHandler).Methods("GET")
	log.Info("Server Start")
	log.Fatal(http.ListenAndServe(":8001", r))
}
