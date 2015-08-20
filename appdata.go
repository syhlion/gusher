package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AppData struct {
	db *sql.DB
}

func NewAppData(db *sql.DB) *AppData {
	return &AppData{db}
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

func (d *AppData) Register(app_name string, request_ip string) (app_key string, err error) {
	cmd := "INSERT INTO appdata(app_name,request_ip,app_key,timestamp,date) VALUES (?,?,?,?,?)"
	tx, err := d.db.Begin()
	if err != nil {
		log.Debug(err)
		return
	}
	stmt, err := tx.Prepare(cmd)
	if err != nil {
		log.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	date := time.Now().Format("2006/01/02 15:04:05")
	timestamp := makeTimestamp()

	seeds := []string{app_name, request_ip, strconv.FormatInt(timestamp, 10), date}
	seed := strings.Join(seeds, ",")
	b := []byte(seed)
	app_key = fmt.Sprintf("%x", md5.Sum(b))

	log.Info(app_key)
	_, err = stmt.Exec(app_name, request_ip, app_key, timestamp, date)
	if err != nil {
		log.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	defer stmt.Close()
	return
}
