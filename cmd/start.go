package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/syhlion/gopusher/core"
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
			logformat.DisableColors = true
			log.Logger.Out = file
		}
	}
	log.Logger.Level = env
	log.Logger.Info("Server Start ", c.String("addr"))
	log.Logger.Fatal(http.ListenAndServe(c.String("addr"), r))
}
