package cmd

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syhlion/gopusher/core"
	"github.com/syhlion/gopusher/model"
	"github.com/syhlion/gopusher/module/config"
	"github.com/syhlion/gopusher/module/log"
	"github.com/syhlion/gopusher/module/requestworker"
	"github.com/syhlion/gopusher/route"
	"net/http"
	"os"
	"time"
)

var CmdStart = cli.Command{
	Name:        "start",
	Usage:       "Start GoPusher server",
	Description: `GoPusher Start`,
	Action:      start,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Value: "./default.json",
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
	create table if not exists appdata (app_name,auth_account,auth_password,connect_hook,request_ip,app_key PRIMARY KEY,timestamp,date)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return
	}
	return
}

//Server執行動作
func start(c *cli.Context) {

	logger := log.Logger
	logformat := &logrus.TextFormatter{FullTimestamp: true}
	logger.Formatter = logformat

	conf := config.GetConfig(c.String("conf"))
	db, err := DBinit(conf.SqlDir)
	if err != nil {
		log.Logger.Fatal(err)
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
		log.Logger.Fatal(err)
	}

	//init env
	env := func() logrus.Level {
		switch conf.Environment {
		case "PRODUCTION":
			return logrus.WarnLevel
			break
		case "DEVELOPMENT":
			return logrus.InfoLevel
			break
		case "DEBUG":
			return logrus.DebugLevel
			break
		}
		return logrus.InfoLevel
	}()

	//init log print
	if conf.LogDir != "console" {
		if file, err := os.OpenFile(conf.LogDir, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0665); err == nil {
			logformat.DisableColors = true
			log.Logger.Out = file
		}
	}
	log.Logger.Level = env
	log.Logger.Info("Server Start ", conf.Listen)

	//server start
	log.Logger.Fatal(http.ListenAndServe(conf.Listen, r))
}
