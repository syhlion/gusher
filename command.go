package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syhlion/requestwork"
	"golang.org/x/crypto/bcrypt"
)

var (
	CmdStart = cli.Command{
		Name:        "start",
		Usage:       "Start GoPusher server",
		Description: `GoPusher Start`,
		Action:      Start,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "conf, c",
				Value: "./config.json",
				Usage: "Input default.json",
			},
		},
	}
	CmdInitConfig = cli.Command{
		Name:        "init",
		Usage:       "Init Gusher Server Config",
		Description: "Init Gusher Server Config",
		Action:      InitConfig,
	}
	Model      *AppData
	ReqWorker  *requestwork.Worker
	GlobalConf *Config
)

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
func Start(c *cli.Context) {

	logformat := &log.TextFormatter{FullTimestamp: true}
	log.SetFormatter(logformat)

	GlobalConf = ConfigGet(c.String("conf"))
	db, err := DBinit(GlobalConf.SqlFile)
	if err != nil {
		log.Fatal(err)
	}

	Model = &AppData{db}
	//init requestwork
	ReqWorker = requestwork.New(&http.Client{
		Timeout: time.Duration(5 * time.Second),
	}, GlobalConf.MaxWaitHook)

	//init router
	publicRouter := PublicRouter()
	privateRouter := PrivateRouter()

	//init env
	env := func() log.Level {
		switch GlobalConf.Environment {
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
	if GlobalConf.LogFile != "console" {
		if file, err := os.OpenFile(GlobalConf.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0665); err == nil {
			logformat.DisableColors = true
			log.SetOutput(file)
		}
	}
	log.SetLevel(env)
	log.Info("Server Start ", GlobalConf.Listen, ", Api Port ", GlobalConf.ApiListen)

	//server start
	srv := &http.Server{
		Addr:         GlobalConf.Listen,
		Handler:      publicRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	privateSrv := &http.Server{
		Addr:         GlobalConf.ApiListen,
		Handler:      privateRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	privateSrvError := make(chan error)
	srvError := make(chan error)
	go func() {
		privateSrvError <- privateSrv.ListenAndServe()
	}()
	go func() {
		srvError <- srv.ListenAndServe()
	}()
	for {
		select {
		case err := <-privateSrvError:
			log.Fatal(err)
			break
		case err := <-srvError:
			log.Fatal(err)
			break
		}
	}
}

func InitConfig(c *cli.Context) {
	conf := Config{}

	//Set Listen Port
	fmt.Print("Please Input Listen Port (Default: ':8001'):")
	fmt.Scanf("%v\n", &conf.Listen)
	if conf.Listen == "" {
		conf.Listen = ":8001"
	}
	fmt.Println("Input: ", conf.Listen)

	//Set Listen Port
	fmt.Print("Please Input Api Listen Port (Default: ':8002'):")
	fmt.Scanf("%v\n", &conf.ApiListen)
	if conf.ApiListen == "" {
		conf.ApiListen = ":8002"
	}
	fmt.Println("Input: ", conf.ApiListen)

	//Set Account
	fmt.Print("Please Input Admin Auth Account (Default: 'account'):")
	fmt.Scanf("%v\n", &conf.AuthAccount)
	if conf.AuthAccount == "" {
		conf.AuthAccount = "account"
	}
	fmt.Println("Input: ", conf.AuthAccount)

	// Set Password (bcrypt encode)
	fmt.Print("Please Input Admin Auth Password (Default: 'password'):")
	fmt.Scanf("%v\n", &conf.AuthPassword)
	if conf.AuthPassword == "" {
		conf.AuthPassword = "password"
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(conf.AuthAccount+conf.AuthPassword), 5)
	if err != nil {
		panic(err)
	}
	conf.AuthPassword = string(hashPassword)
	fmt.Println("Input: ", "*****************")

	//Set Sqlite file
	fmt.Print("Please Input SQL File Location (Default: './appdata.sqlite'):")
	fmt.Scanf("%v\n", &conf.SqlFile)
	if conf.SqlFile == "" {
		conf.SqlFile = "./appdata.sqlite"
	}
	fmt.Println("Input: ", conf.SqlFile)

	// Set log file or console
	fmt.Print("Please Input Log File Location OR Console Log (Default: 'console' Option: console || ./gusher.log):")
	fmt.Scanf("%v\n", &conf.LogFile)
	if conf.LogFile == "" {
		conf.LogFile = "console"
	}
	fmt.Println("Input: ", conf.LogFile)

	//Set env
	fmt.Print("Please Input Environment (Default: 'DEBUG' Option: DEBUG || DEVELOPMENT || PRODUCATION):")
	fmt.Scanf("%v\n", &conf.Environment)
	if conf.Environment == "" {
		conf.Environment = "DEBUG"
	}

	if !(conf.Environment == "DEBUG" || conf.Environment == "DEVELOPMENT" || conf.Environment == "PRODUCATION") {
		conf.Environment = "DEBUG"
	}
	fmt.Println("Input: ", conf.Environment)

	//Set WebHook access Resource
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
	err = ConfigWrite(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scuess")

}
