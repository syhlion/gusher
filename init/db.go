package init

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func DBinit() (err error) {
	db, err := sql.Open("sqlite3", "./appdata.sqlite")
	if err != nil {
		return
	}

	sqlStmt := `
	create table if not exists appdata (app_name,auth_account,auth_password,request_ip,app_key PRIMARY KEY,timestamp,date)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return
	}
	return
}
