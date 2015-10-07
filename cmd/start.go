package cmd

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syhlion/gusher/core"
	"github.com/syhlion/gusher/model"
	"github.com/syhlion/gusher/module/config"
	"github.com/syhlion/gusher/module/requestworker"
	"github.com/syhlion/gusher/route"
)

var CmdStart = cli.Command{
	Name:        "start",
	Usage:       "Start GoPusher server",
	Description: `GoPusher Start`,
	Action:      start,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Value: "./config.json",
			Usage: "Input default.json",
		},
	},
}

func DBinit(sqlDir string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", sqlDir)
	if err != nil {
		return
	}

	sqlStmt := `
	create table if not exists appdata (app_key PRIMARY KEY,auth_account,auth_password,connect_hook,request_ip,timestamp,date)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return
	}
	return
}

//Server執行動作
func start(c *cli.Context) {

	logformat := &log.TextFormatter{FullTimestamp: true}
	log.SetFormatter(logformat)

	conf := config.Get(c.String("conf"))
	db, err := DBinit(conf.SqlFile)
	if err != nil {
		log.Fatal(err)
	}

	//init core
	collection := core.NewCollection()
	go collection.Run()

	//init model appdata
	appdata := model.NewAppData(db)

	//init requestworker
	worker := &requestworker.Worker{
		Threads:  conf.MaxWaitHook,
		JobQuene: make(chan *requestworker.Job, 1024),
		HttpClient: &http.Client{
			Timeout: time.Duration(5 * time.Second),
		},
	}

	//work start wait
	go worker.Start()

	//init router
	r := route.Router(appdata, collection, conf, worker)
	if err != nil {
		log.Fatal(err)
	}

	//init env
	env := func() log.Level {
		switch conf.Environment {
		case "PRODUCTION":
			return log.WarnLevel
			break
		case "DEVELOPMENT":
			return log.InfoLevel
			break
		case "DEBUG":
			return log.DebugLevel
			break
		}
		return log.InfoLevel
	}()

	//init log print
	if conf.LogFile != "console" {
		if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0665); err == nil {
			logformat.DisableColors = true
			log.SetOutput(file)
		}
	}
	log.SetLevel(env)
	log.Info("Server Start ", conf.Listen)

	//server start
	srv := &http.Server{
		Addr:         conf.Listen,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
