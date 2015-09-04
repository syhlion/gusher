package cmd

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/syhlion/gusher/module/config"
	"strconv"
)

var InitStart = cli.Command{
	Name:        "init",
	Usage:       "Init Gusher Server Config",
	Description: "Init Gusher Server Config",
	Action:      initStart,
}

func initStart(c *cli.Context) {
	conf := config.Config{}
auth_account:
	fmt.Print("Please Input Admin Auth Account:")
	fmt.Scan(&conf.AuthAccount)
	if conf.AuthAccount == "" {
		goto auth_account
	}

auth_password:
	fmt.Print("Please Input Admin Auth Password:")
	fmt.Scan(&conf.AuthPassword)
	if conf.AuthPassword == "" {
		goto auth_password
	}

sql_file:
	fmt.Print("Please Input SQL File Location (ex: ./appdata.sqlite):")
	fmt.Scan(&conf.SqlFile)
	if conf.SqlFile == "" {
		goto sql_file
	}

log_file:
	fmt.Print("Please Input Log File Location OR Console Log (ex: console || ./gusher.log):")
	fmt.Scan(&conf.LogFile)
	if conf.LogFile == "" {
		goto log_file
	}

env:
	fmt.Print("Please Input Environment (DEBUG || DEVELOPMENT || PRODUCATION):")
	fmt.Scan(&conf.Environment)
	if conf.Environment == "" {
		goto env
	}

	if !(conf.Environment == "DEBUG" || conf.Environment == "DEVELOPMENT" || conf.Environment == "PRODUCATION") {
		goto env
	}

max_wait_hook:
	fmt.Print("Please Input the Nnumber WEB HOOK Request BOT (ex: 100):")
	var num string
	fmt.Scan(&num)
	n, err := strconv.Atoi(num)
	if err != nil {
		goto max_wait_hook
	}
	conf.MaxWaitHook = n

	fmt.Printf("%+v\n", conf)
	err = config.Write(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scuess")

}
