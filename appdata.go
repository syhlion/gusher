package main

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syhlion/go-common"
)

var (
	appData *AppData = nil
)

type AppDataResult struct {
	AppKey       string `json:"app_key"`
	AuthAccount  string `json:"auth_account"`
	AuthPassword string `json:"auth_password"`
	ConnectHook  string `json:"connect_hook"`
	RequestIP    string `json:"request_ip"`
	Date         string `json:"date"`
	Timestamp    string `json:"timestamp"`
}
type AppData struct {
	db *sql.DB
}

func (d *AppData) IsExist(app_key string) bool {
	sql := "SELECT EXISTS(SELECT 1  FROM  `appdata` WHERE `app_key`= $1)"
	var result int
	err := d.db.QueryRow(sql, app_key).Scan(&result)
	if err != nil {
		log.Debug(app_key, " ", err)
		return false
	}

	if result == 0 {
		log.Debug(app_key, " no exist")
		return false
	}
	return true

}

func (d *AppData) Delete(app_key string) (err error) {

	sql := "DELETE FROM `appdata` where app_key = ?"

	tx, err := d.db.Begin()
	if err != nil {
		log.Debug(err)
		return
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Debug(app_key, " ", err)
		return
	}

	_, err = stmt.Exec(app_key)
	if err != nil {
		log.Debug(app_key, " ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Debug(app_key, " ", err)
		return
	}
	return

}

func (d *AppData) Get(app_key string) (r AppDataResult, err error) {

	sql := "SELECT * FROM `appdata` WHERE app_key = ?"
	stmt, err := d.db.Prepare(sql)
	if err != nil {
		log.Debug(err)
		return
	}

	rows, err := stmt.Query(app_key)
	if err != nil {
		log.Debug(err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&r.AppKey, &r.AuthAccount, &r.AuthPassword, &r.ConnectHook, &r.RequestIP, &r.Timestamp, &r.Date)
		if err != nil {
			log.Debug(err)
			return
		}
	}
	return
}

func (d *AppData) GetAll() (r []AppDataResult, err error) {

	sql := "SELECT * FROM `appdata`"
	rows, err := d.db.Query(sql)
	if err != nil {
		log.Debug(err)
		return
	}
	var apps AppDataResult
	for rows.Next() {
		err = rows.Scan(&apps.AppKey, &apps.AuthAccount, &apps.AuthPassword, &apps.ConnectHook, &apps.RequestIP, &apps.Timestamp, &apps.Date)
		if err != nil {
			log.Debug(err)
			return
		}
		r = append(r, apps)
	}
	return

}

func (d *AppData) Register(app_key string, auth_account string, auth_password string, connect_hook string, request_ip string) (err error) {
	cmd := "INSERT INTO appdata(app_key,auth_account,auth_password,connect_hook,request_ip,timestamp,date) VALUES (?,?,?,?,?,?,?)"
	tx, err := d.db.Begin()
	if err != nil {
		log.Debug(err)
		return
	}
	stmt, err := tx.Prepare(cmd)
	if err != nil {
		log.Debug(app_key, " ", request_ip, " ", err)
		return
	}
	date := time.Now().Format("2006/01/02 15:04:05")

	_, err = stmt.Exec(app_key, auth_account, auth_password, connect_hook, request_ip, common.Time(), date)
	if err != nil {
		log.Debug(app_key, " ", request_ip, " ", err)
		tx.Rollback()
		stmt.Close()
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Debug(app_key, " ", request_ip, " ", err)
		tx.Rollback()
		stmt.Close()
		return
	}
	defer stmt.Close()
	return
}
