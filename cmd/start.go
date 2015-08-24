package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/syhlion/gopusher/core"
	"github.com/syhlion/gopusher/init"
	"github.com/syhlion/gopusher/module/log"
	"github.com/syhlion/gopusher/router"
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

//Server執行動作
func start(c *cli.Context) {

	collection := core.Collection
	logger := log.Logger
	logformat := &logrus.TextFormatter{FullTimestamp: true}
	logger.Formatter = logformat

	go collection.Run()
	r := router.Router()
	conf := init.GetConfig(c.String("conf"))
	err := init.DBinit()
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
		if file, err := os.OpenFile(c.String("log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0665); err == nil {
			logformat.DisableColors = true
			log.Logger.Out = file
		}
	}
	log.Logger.Level = env
	log.Logger.Info("Server Start ", conf.Listen)
	log.Logger.Fatal(http.ListenAndServe(conf.Listen, r))
}
