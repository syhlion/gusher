package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"net/http"
	"os"
)

const (
	APP_VER = "0.2.1"
)

// 管理每個 app的集合初始化
var collection = NewCollection()

var logformat = &logrus.TextFormatter{FullTimestamp: true}

var log = logrus.New()

var appdata *AppData

//初始化執行動作
func Start(c *cli.Context) {

	log.Formatter = logformat
	go collection.run()
	db, err := DBinit()
	if err != nil {
		log.Fatal(err)
	}
	appdata = NewAppData(db)
	r := Router()
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

//進入點
func main() {

	gusher := cli.NewApp()
	gusher.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr, a",
			Value: ":8001",
			Usage: "Input like 127.0.0.1:8001 or :8011",
		},
		cli.StringFlag{
			Name:  "env, e",
			Value: "PRODUCTION",
			Usage: "PRODUCTION | DEVELOPMENT | DEBUG",
		},
		cli.StringFlag{
			Name:  "log, l",
			Value: "console",
			Usage: "Input like /home/user/gusher.log | console",
		},
	}
	gusher.Name = "gusher"
	gusher.Version = APP_VER
	gusher.Action = Start

	gusher.Run(os.Args)

}
