package user

import (
	"database/sql"
	"fmt"
	"track/dbs"
)
var db *sql.DB
var Db dbs.Database
func init(){
	var err error
	db, err = sql.Open("mysql", "root:09070810.@tcp(localhost:3306)/track")
	if err != nil {
		fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
	}
	Db.MysqlDb=db
}
