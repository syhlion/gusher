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
	fmt.Print("Please Input Listen Port (Default: ':8001'):")
	fmt.Scanf("%v\n", &conf.Listen)
	if conf.Listen == "" {
		conf.Listen = ":8001"
	}
	fmt.Println("Input: ", conf.Listen)

	fmt.Print("Please Input Admin Auth Account (Default: 'account'):")
	fmt.Scanf("%v\n", &conf.AuthAccount)
	if conf.AuthAccount == "" {
		conf.AuthAccount = "account"
	}
	fmt.Println("Input: ", conf.AuthAccount)

	fmt.Print("Please Input Admin Auth Password (Default: 'password'):")
	fmt.Scanf("%v\n", &conf.AuthPassword)
	if conf.AuthPassword == "" {
		conf.AuthPassword = "password"
	}
	fmt.Println("Input: ", conf.AuthPassword)

	fmt.Print("Please Input SQL File Location (Default: './appdata.sqlite'):")
	fmt.Scanf("%v\n", &conf.SqlFile)
	if conf.SqlFile == "" {
		conf.SqlFile = "./appdata.sqlite"
	}
	fmt.Println("Input: ", conf.SqlFile)

	fmt.Print("Please Input Log File Location OR Console Log (Default: 'console' Option: console || ./gusher.log):")
	fmt.Scanf("%v\n", &conf.LogFile)
	if conf.LogFile == "" {
		conf.LogFile = "console"
	}
	fmt.Println("Input: ", conf.LogFile)

	fmt.Print("Please Input Environment (Default: 'DEBUG' Option: DEBUG || DEVELOPMENT || PRODUCATION):")
	fmt.Scanf("%v\n", &conf.Environment)
	if conf.Environment == "" {
		conf.Environment = "DEBUG"
	}

	if !(conf.Environment == "DEBUG" || conf.Environment == "DEVELOPMENT" || conf.Environment == "PRODUCATION") {
		conf.Environment = "DEBUG"
	}
	fmt.Println("Input: ", conf.Environment)

	fmt.Print("Please Input the Nnumber WEB HOOK Request BOT (Default: 100):")
	var num string
	fmt.Scanf("%v\n", &num)
	n, err := strconv.Atoi(num)
	if err != nil {
		conf.MaxWaitHook = 100
	} else {
		conf.MaxWaitHook = n
	}
	fmt.Println("Input: ", conf.MaxWaitHook)

	fmt.Printf("Result Json: %+v\n", conf)
	err = config.Write(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scuess")

}
