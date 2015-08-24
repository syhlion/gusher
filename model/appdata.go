package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syhlion/go-common"
	"github.com/syhlion/gopusher/module/log"
	"strings"
	"time"
)

var (
	AppData *appData = nil
)

func init() {
	if AppData == nil {

		db, err := sql.Open("sqlite3", "./appdata.sqlite")
		if err != nil {
			log.Logger.Error(err)
			return
		}
		AppData = newAppData(db)
	}
}

type AppDataResult struct {
	AppKey    string `json:"app_key"`
	AppName   string `json:"app_name"`
	RequestIP string `json:"request_ip"`
	Date      string `json:"date"`
	Timestamp string `json:"timestamp"`
}
type appData struct {
	db *sql.DB
}

func newAppData(db *sql.DB) *appData {
	return &appData{db}
}

func (d *appData) IsExist(app_key string) bool {
	sql := "SELECT EXISTS(SELECT 1  FROM  `appdata` WHERE `app_key`= $1)"
	var result int
	err := d.db.QueryRow(sql, app_key).Scan(&result)
	if err != nil {
		log.Logger.Debug(app_key, " ", err)
		return false
	}

	if result == 0 {
		log.Logger.Debug(app_key, " no exist")
		return false
	}
	return true

}

func (d *appData) Delete(app_key string) (err error) {
	sql := "DELETE FROM `appdata` where app_key = ?"

	tx, err := d.db.Begin()
	if err != nil {
		log.Logger.Debug(err)
		return
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Logger.Debug(app_key, " ", err)
		return
	}

	_, err = stmt.Exec(app_key)
	if err != nil {
		log.Logger.Debug(app_key, " ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Logger.Debug(app_key, " ", err)
		return
	}
	return

}

func (d *appData) GetAll() (r []AppDataResult, err error) {

	sql := "SELECT * FROM `appdata`"
	rows, err := d.db.Query(sql)
	if err != nil {
		log.Logger.Debug(err)
		return
	}
	var apps AppDataResult
	for rows.Next() {
		err = rows.Scan(&apps.AppName, &apps.RequestIP, &apps.AppKey, &apps.Timestamp, &apps.Date)
		if err != nil {
			log.Logger.Debug(err)
			return
		}
		r = append(r, apps)
	}
	return

}

func (d *appData) Register(app_name string, request_ip string) (app_key string, err error) {
	cmd := "INSERT INTO appdata(app_name,request_ip,app_key,timestamp,date) VALUES (?,?,?,?,?)"
	tx, err := d.db.Begin()
	if err != nil {
		log.Logger.Debug(err)
		return
	}
	stmt, err := tx.Prepare(cmd)
	if err != nil {
		log.Logger.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	date := time.Now().Format("2006/01/02 15:04:05")

	seeds := []string{app_name, request_ip, common.TimeToString(), date}
	seed := strings.Join(seeds, ",")
	app_key = common.EncodeMd5(seed)

	_, err = stmt.Exec(app_name, request_ip, app_key, common.Time(), date)
	if err != nil {
		log.Logger.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Logger.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	defer stmt.Close()
	return
}
