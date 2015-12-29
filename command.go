package main

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syhlion/requestwork"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	CmdStart = cli.Command{
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
	InitStart = cli.Command{
		Name:        "init",
		Usage:       "Init Gusher Server Config",
		Description: "Init Gusher Server Config",
		Action:      initStart,
	}
	Model     *AppData            = nil
	ReqWorker *requestwork.Worker = nil
	Conf      *Config             = nil
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
func start(c *cli.Context) {

	logformat := &log.TextFormatter{FullTimestamp: true}
	log.SetFormatter(logformat)

	Conf := ConfigGet(c.String("conf"))
	db, err := DBinit(Conf.SqlFile)
	if err != nil {
		log.Fatal(err)
	}

	Model = &AppData{db}
	//init requestwork
	ReqWorker = requestwork.New(&http.Client{
		Timeout: time.Duration(5 * time.Second),
	}, Conf.MaxWaitHook)

	//init router
	r := Router()

	//init env
	env := func() log.Level {
		switch Conf.Environment {
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
	if Conf.LogFile != "console" {
		if file, err := os.OpenFile(Conf.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0665); err == nil {
			logformat.DisableColors = true
			log.SetOutput(file)
		}
	}
	log.SetLevel(env)
	log.Info("Server Start ", Conf.Listen)

	//server start
	srv := &http.Server{
		Addr:         Conf.Listen,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func initStart(c *cli.Context) {
	conf := Config{}

	//Set Listen Port
	fmt.Print("Please Input Listen Port (Default: ':8001'):")
	fmt.Scanf("%v\n", &conf.Listen)
	if conf.Listen == "" {
		conf.Listen = ":8001"
	}
	fmt.Println("Input: ", conf.Listen)

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

	//Set AllowAccessApiIP
	fmt.Print("Please Input Allow Access Api IP(Default: '' <- it means allow all Ex: 192.168  or 127.0.0.1 :")
	var ip string
	var ips []string
	fmt.Scanf("%v\n", &ip)
	if ip == "" {
		ip = ""
	}
	ips = append(ips, ip)
	conf.AllowAccessApiIP = ips

	fmt.Println("Input: ", ips)
	fmt.Printf("Result Json: %+v\n", conf)
	err = ConfigWrite(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scuess")

}
