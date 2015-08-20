package main

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"time"
)

const (
	APP_VER = "0.2.0"
)

var collection = NewCollection()

var log = logrus.New()

var appdata *AppData

func makeTimestamp() (t int64) {
	t = time.Now().UnixNano() / int64(time.Millisecond)
	return
}
func main() {
	logformat := &logrus.TextFormatter{FullTimestamp: true}
	log.Formatter = logformat
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
	r.HandleFunc("/{app_key}/listonlineuser", ListClientHandler).Methods("GET")

	gusher := cli.NewApp()
	gusher.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: ":8001",
			Usage: "Input like 127.0.0.1:8001 or :8011",
		},
		cli.StringFlag{
			Name:  "env",
			Value: "PRODUCTION",
			Usage: "PRODUCTION | DEVELOPMENT | DEBUG",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "console",
			Usage: "Input like /home/user/gusher.log | console",
		},
	}
	gusher.Name = "gusher"
	gusher.Version = APP_VER
	gusher.Action = func(c *cli.Context) {
		env := func() logrus.Level {
			switch c.String("env") {
			case "PRODUCTION":
				return logrus.InfoLevel
				break
			case "DEVELOPMENT":
				return logrus.InfoLevel
			case "DEBUG":
				return logrus.DebugLevel
			}
			return logrus.WarnLevel
		}()
		if c.String("log") != "console" {
			if file, err := os.OpenFile(c.String("log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0665); err == nil {
				log.Out = file
				logformat.DisableColors = true
			}
		}
		log.Level = env
		log.Info("Server Start ", c.String("addr"))
		log.Fatal(http.ListenAndServe(c.String("addr"), r))
	}

	gusher.Run(os.Args)

}
