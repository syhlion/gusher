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
	"github.com/syhlion/gopusher/route"
	"net/http"
	"os"
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
	create table if not exists appdata (app_name,auth_account,auth_password,request_ip,app_key PRIMARY KEY,timestamp,date)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return
	}
	return
}

//Server執行動作
func start(c *cli.Context) {

	collection := core.NewCollection()
	logger := log.Logger
	logformat := &logrus.TextFormatter{FullTimestamp: true}
	logger.Formatter = logformat

	go collection.Run()
	conf := config.GetConfig(c.String("conf"))
	db, err := DBinit(conf.SqlDir)
	if err != nil {
		log.Logger.Fatal(err)
	}
	appdata := model.NewAppData(db)
	r := route.Router(appdata, collection, conf)
	if err != nil {
		log.Logger.Fatal(err)
	}
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
	if conf.LogDir != "console" {
		if file, err := os.OpenFile(conf.LogDir, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0665); err == nil {
			logformat.DisableColors = true
			log.Logger.Out = file
		}
	}
	log.Logger.Level = env
	log.Logger.Info("Server Start ", conf.Listen)
	log.Logger.Fatal(http.ListenAndServe(conf.Listen, r))
}
